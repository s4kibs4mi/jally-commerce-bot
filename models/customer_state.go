package models

type CustomerState string

const (
	CustomerStateStart          CustomerState = "start"
	CustomerStateChooseMenu     CustomerState = "choose_menu"
	CustomerStateCheckout       CustomerState = "checkout"
	CustomerStateAdditionalData CustomerState = "additional_data"
	CustomerStatePay            CustomerState = "pay"
	CustomerStateOrderConfirmed CustomerState = "order_confirmed"
	CustomerStateSearchProducts CustomerState = "search_products"
)
