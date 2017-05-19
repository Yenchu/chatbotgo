package main

import (
	"flag"

	"github.com/yenchu/chatbotgo/apiai"
	"github.com/yenchu/chatbotgo/fb"
)

var (
	port             = flag.Int("port", 8080, "The server port")
	validationToken  = flag.String("validation_token", "", "The token used to valify webhook")
	pageAccessToken  = flag.String("page_access_token", "", "The token used to access Facebook API")
	apiaiAccessToken = flag.String("apiap_access_token", "", "The token used to access API.AI")
)

func main() {
	flag.Parse()
	startFBMessenger()
}

func startFBMessenger() {
	aiOptions := apiai.Options{AccessToken: *apiaiAccessToken, Lang: "zh-TW"}
	apiaiClient := apiai.NewClient(aiOptions)

	handleTextFunc := func(recipientID string, messageText string) string {
		qryRes, err := apiaiClient.Query(recipientID, messageText)
		if err != nil {
			return "Sorry! I cannot understand what you said."
		}

		respText := qryRes.Result.Fulfillment.Speech
		return respText
	}

	fbOptions := fb.Options{Port: *port, ValidationToken: *validationToken, PageAccessToken: *pageAccessToken,
		HandleText: handleTextFunc}
	messanger := fb.NewMessenger(fbOptions)
	messanger.Start()
}

// for test
func startApiAIClient() {
	aiOptions := apiai.Options{AccessToken: *apiaiAccessToken, Lang: "zh-TW"}
	apiaiClient := apiai.NewClient(aiOptions)
	apiaiClient.Query("andrew", "weather in Taipei")
}
