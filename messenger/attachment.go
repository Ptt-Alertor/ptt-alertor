package messenger

type Attachment struct {
	Type    string `json:"type"`
	Payload `json:"payload"`
}

type Payload struct {
	TemplateType    string    `json:"template_type"`
	TopElementStyle string    `json:"top_element_style,omitempty"`
	Elements        []Element `json:"elements"`
}

type Element struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
}
