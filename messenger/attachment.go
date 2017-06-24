package messenger

type Attachment struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

type ListPayload struct {
	TemplateType    string    `json:"template_type"`
	TopElementStyle string    `json:"top_element_style,omitempty"`
	Elements        []Element `json:"elements"`
}

type Element struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
}

type ButtonPayload struct {
	TemplateType string `json:"template_type"`
	Text         string `json:"text"`
	Buttons      `json:"buttons"`
}

type Buttons []Button

type Button struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Payload string `json:"payload"`
}
