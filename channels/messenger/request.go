package messenger

type Request struct {
	Recipient   `json:"recipient"`
	Message     `json:"message"`
	MessageType string `json:"message_type"`
	Tag         string `json:"tag"`
}

type Recipient struct {
	ID string `json:"id"`
}
