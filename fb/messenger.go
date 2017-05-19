package fb

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Messenger interface {
	Start()
}

type Options struct {
	Port            int
	Version         string
	ValidationToken string
	PageAccessToken string
	HandleText      HandleTextFunc
}

type HandleTextFunc func(recipientID string, messageText string) string

func NewMessenger(options Options) Messenger {
	port := "8080"
	if options.Port > 0 {
		port = strconv.Itoa(options.Port)
	}

	version := "v2.6"
	if len(options.Version) > 0 {
		version = options.Version
	}
	facebookMessagesURL := "https://graph.facebook.com/" + version + "/me" + "/messages?access_token=" + options.PageAccessToken

	var client = &http.Client{Timeout: time.Second * 20}
	m := messenger{client: client, port: port, facebookMessagesURL: facebookMessagesURL,
		validationToken: options.ValidationToken, pageAccessToken: options.PageAccessToken,
		handleTextFunc: options.HandleText}
	return &m
}

type messenger struct {
	client              *http.Client
	port                string
	facebookMessagesURL string
	validationToken     string
	pageAccessToken     string
	handleTextFunc      HandleTextFunc
}

func (m *messenger) Start() {
	http.HandleFunc("/webhook", m.webhookHandler)

	srv := &http.Server{
		Addr: ":" + m.port,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Printf("An error occured starting Facebook messanger at port %v", m.port)
		log.Printf("Error: %v", err.Error())
	}
}

func (m *messenger) webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		m.handleValidation(w, r)
	} else if r.Method == "POST" {
		m.handleSubscription(w, r)
	} else {
		log.Printf("Unsupported HTTP method %v", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (m *messenger) handleValidation(w http.ResponseWriter, r *http.Request) {
	if queryString(r, "hub.mode") == "subscribe" && queryString(r, "hub.verify_token") == m.validationToken {
		log.Println("Validating webhook")
		io.WriteString(w, queryString(r, "hub.challenge"))
		return
	}

	log.Println("Failed validation. Make sure the validation tokens match.")
	w.WriteHeader(http.StatusForbidden)
}

func (m *messenger) handleSubscription(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read the request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var recvMessage ReceivedMessage
	if err := json.Unmarshal(body, &recvMessage); err != nil {
		log.Printf("Could not decode request body as JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure this is a page subscription
	if recvMessage.Object == "page" {
		// Iterate over each entry - there may be multiple if batched
		for _, entry := range recvMessage.Entry {
			/*pageID := entry.ID
			timeOfEvent := entry.Time
			log.Printf("Receive message for page %v at %v", pageID, timeOfEvent)*/

			// Iterate over each messaging event
			for _, event := range entry.Messaging {
				// There may have 10 relevant webhook events, only handling Message Received events for demo
				if event.Message != nil {
					senderID, messageText := m.receiveMessage(event)

					respText := m.handleTextFunc(senderID, messageText)

					m.sendTextMessage(w, senderID, respText)
				} else {
					log.Printf("Webhook received unknown event: %v", event)
				}
			}
		}
	} else {
		log.Printf("Webhook received non-page message: %v", recvMessage)
	}
}

func (m *messenger) receiveMessage(event *Messaging) (string, string) {
	senderID := event.Sender.ID
	message := event.Message
	log.Printf("Received message from %v: %v", senderID, message.Text)

	messageText := message.Text
	attachments := message.Attachments
	if messageText != "" {
		return senderID, messageText
	} else if len(attachments) > 0 {
		return senderID, "Message with attachment received"
	}
	return senderID, "Unrecognised message received"
}

func (m *messenger) sendTextMessage(w http.ResponseWriter, recipientID string, messageText string) {
	log.Printf("Send message to %v: %v", recipientID, messageText)
	sentMessage := SentTextMessage{Recipient: &Recipient{recipientID}, Message: &TextMessage{Text: messageText}}
	m.sendMessage(w, sentMessage)
}

func (m *messenger) sendMessage(w http.ResponseWriter, message SentTextMessage) {
	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Could not encode message as JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := m.client.Post(m.facebookMessagesURL, "application/json; charset=UTF-8", bytes.NewReader(body))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
	defer res.Body.Close()
}

func queryString(r *http.Request, param string) string {
	return r.URL.Query().Get(param)
}
