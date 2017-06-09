package messenger

type Request struct {
	Recipient `json:"recipient"`
	Message   `json:"message"`
}

type Recipient struct {
	ID string `json:"id"`
}
