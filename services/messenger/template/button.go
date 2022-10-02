package template

const TemplateTypeButton TemplateType = "button"

// ButtonTemplate is a template
type ButtonTemplate struct {
	Text         string       `json:"text,omitempty"`
	Buttons      []Button     `json:"buttons,omitempty"`
	TemplateType TemplateType `json:"template_type"`
}

func (ButtonTemplate) Type() TemplateType {
	return TemplateTypeButton
}

func (ButtonTemplate) SupportsButtons() bool {
	return true
}
