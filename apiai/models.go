package apiai

type Button struct {
	Text     string `json:"text"`
	Postback string `json:"postback"`
}

type Context struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
	Lifespan   int32             `json:"lifespan"`
}

type Entity struct {
	Name    string         `json:"name"`
	Entries []*EntityEntry `json:"entries"`
	Extend  bool           `json:"extend"`
	IsEnum  bool           `json:"isEnum"`
}

type EntityEntry struct {
	Value    string   `json:"value"`
	Synonyms []string `json:"synonyms"`
}

type Event struct {
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

type Fulfillment struct {
	Speech   string     `json:"speech"`
	Messages []*Message `json:"messages"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Message struct {
	Type int32 `json:"type"`
}

// TextMessage is a Message with type=0
type TextMessage struct {
	Message
	Speech string `json:"speech"`
}

// ImageMessage is a Message with type=3
type ImageMessage struct {
	Message
	ImageURL string `json:"imageUrl"`
}

// CardMessage is a Message with type=1
type CardMessage struct {
	Message
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Buttons  []*Button `json:"buttons"`
}

// QuickRepliesMessage is a Message with type=2
type QuickRepliesMessage struct {
	Message
	Title   string   `json:"title"`
	Replies []string `json:"replies"`
}

// CustomPayloadMessage is a Message with type=4
type CustomPayloadMessage struct {
	Message
	Payload string `json:"payload"` // JSON string
}

type Metadata struct {
	IntentID                  string `json:"intentId"`
	IntentName                string `json:"intentName"`
	WebhookUsed               string `json:"webhookUsed"`
	WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
}

type OriginalRequest struct {
	Source string            `json:"source"`
	Data   map[string]string `json:"data"`
}

type QuestionMetadata struct {
	SessionID string    `json:"sessionId"`
	Lang      string    `json:"lang"`
	Entities  []*Entity `json:"entities"`
	Timezone  string    `json:"timezone"`
	Location  *Location `json:"location"`
}

type QueryRequest struct {
	Query           []string         `json:"query,omitempty"`
	Event           *Event           `json:"event,omitempty"`
	SessionID       string           `json:"sessionId,omitempty"`
	Lang            string           `json:"lang,omitempty"`
	Contexts        []*Context       `json:"contexts,omitempty"`
	ResetContexts   bool             `json:"resetContexts,omitempty"`
	Entities        []*Entity        `json:"entities,omitempty"`
	Timezone        string           `json:"timezone,omitempty"`
	Location        *Location        `json:"location,omitempty"`
	OriginalRequest *OriginalRequest `json:"originalRequest,omitempty"`
}

type QueryResponse struct {
	ID        string  `json:"id"`
	Timezone  string  `json:"timezone"`
	Lang      string  `json:"lang"`
	Result    *Result `json:"result"`
	Status    *Status `json:"status"`
	SessionID string  `json:"sessionId"`
}

type Result struct {
	Source           string            `json:"source"`
	ResolvedQuery    string            `json:"resolvedQuery"`
	Action           string            `json:"action"`
	ActionIncomplete bool              `json:"actionIncomplete"`
	Parameters       map[string]string `json:"parameters"`
	Contexts         []*Context        `json:"contexts"`
	Fulfillment      *Fulfillment      `json:"fulfillment"`
	Score            float32           `json:"score"`
	Metadata         *Metadata         `json:"metadata"`
}

type Status struct {
	Code         int32  `json:"code"`
	ErrorType    string `json:"errorType"`
	ErrorDetails string `json:"errorDetails"`
	ErrorID      string `json:"errorID"`
}
