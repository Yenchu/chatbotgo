package fb

type ReceivedMessage struct {
	Object string   `json:"object"`
	Entry  []*Entry `json:"entry"`
}

type Entry struct {
	ID        string       `json:"id"`
	Time      int64        `json:"time"`
	Messaging []*Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    *Sender    `json:"sender"`
	Recipient *Recipient `json:"recipient"`
	Timestamp int64      `json:"timestamp"`
	Message   *Message   `json:"message"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	MID         string        `json:"mid"`
	Seq         int64         `json:"seq"`
	Text        string        `json:"text"`
	Attachments []*Attachment `json:"attachments"`
	QuickReply  *QuickReply   `json:"quick_reply"`
}

type Attachment struct {
	Type    string `json:"type"`
	Payload struct {
		// for image, audio, video or file payload
		URL string `json:"url"`

		// for location payload
		Latitude  float64 `json:"coordinates.lat"`
		Longitude float64 `json:"coordinates.long"`
	} `json:"payload"`
}

type QuickReply struct {
	Payload string `json:"payload"`
}

type SentTextMessage struct {
	Recipient *Recipient   `json:"recipient"`
	Message   *TextMessage `json:"message"`
}

type TextMessage struct {
	Text     string `json:"text"`
	Metadata string `json:"metadata"`
}
