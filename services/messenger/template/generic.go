package template

const TemplateTypeGeneric TemplateType = "generic"
const TemplateTypeReceipt TemplateType = "receipt"

type GenericTemplate struct {
	// Title is limited to 45 characters
	Title                string                `json:"title"`
	ItemURL              string                `json:"item_url,omitempty"`
	ImageURL             string                `json:"image_url,omitempty"`
	DefaultActionGeneric *DefaultActionGeneric `json:"default_action,omitempty"`
	// Subtitle is limited to 80 characters
	Subtitle string   `json:"subtitle,omitempty"`
	Buttons  []Button `json:"buttons,omitempty"`
}

type DefaultActionGeneric struct {
	Type ButtonType `json:"type,omitempty"` // Must be "web_url"
	URL  string     `json:"url,omitempty"`
}

func (GenericTemplate) Type() TemplateType {
	return TemplateTypeGeneric
}

func (GenericTemplate) SupportsButtons() bool {
	return true
}

type ReceiptTemplate struct {
	TemplateType  TemplateType              `json:"template_type,omitempty"`
	RecipientName string                    `json:"recipient_name,omitempty"`
	OrderNumber   string                    `json:"order_number,omitempty"`
	Currency      string                    `json:"currency,omitempty"`
	PaymentMethod string                    `json:"payment_method,omitempty"`
	Timestamp     int64                     `json:"timestamp,omitempty"`
	Address       *OrderAddressTemplate     `json:"address,omitempty"`
	Summary       *OrderSummaryTemplate     `json:"summary,omitempty"`
	Elements      []OrderElementTemplate    `json:"elements,omitempty"`
	Adjustments   []OrderAdjustmentTemplate `json:"adjustments,omitempty"`
}

type OrderAddressTemplate struct {
	Street1    string `json:"street_1,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
}

type OrderSummaryTemplate struct {
	Subtotal     float64 `json:"subtotal,omitempty"`
	ShippingCost float64 `json:"shipping_cost,omitempty"`
	TotalTax     float64 `json:"total_tax,omitempty"`
	TotalCost    float64 `json:"total_cost,omitempty"`
}

type OrderElementTemplate struct {
	Title    string  `json:"title,omitempty"`
	Subtitle string  `json:"subtitle,omitempty"`
	Quantity int64   `json:"quantity,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Currency string  `json:"currency,omitempty"`
	ImageUrl string  `json:"image_url,omitempty"`
}

type OrderAdjustmentTemplate struct {
	Name   string  `json:"name,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}

func (ReceiptTemplate) Type() TemplateType {
	return TemplateTypeReceipt
}

func (ReceiptTemplate) SupportsButtons() bool {
	return false
}
