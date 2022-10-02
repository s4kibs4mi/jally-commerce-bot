package models

type Search struct {
	Query   string              `json:"query"`
	Filters []map[string]string `json:"filters"`
}

func (Search) GetGraphQLType() string {
	return "Search"
}

type Sort struct {
	By        string `json:"by"`
	Direction string `json:"direction"`
}

func (Sort) GetGraphQLType() string {
	return "Sort"
}

type Pagination struct {
	PerPage int `json:"perPage"`
	Page    int `json:"page"`
}

func (Pagination) GetGraphQLType() string {
	return "Pagination"
}

type NewCartParams struct {
	CartItems []CartItem `json:"cartItems"`
}

func (NewCartParams) GetGraphQLType() string {
	return "NewCartParams"
}

type UpdateCartParams struct {
	CartItems []CartItem `json:"cartItems"`
}

func (UpdateCartParams) GetGraphQLType() string {
	return "UpdateCartParams"
}

type CartItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type QueryCartItem struct {
	Product       Product `json:"product"`
	Quantity      int     `json:"quantity"`
	PurchasePrice int64   `json:"purchasePrice"`
}

func (CartItem) GetGraphQLType() string {
	return "CartItemParams"
}

type GuestCheckoutPlaceOrderParams struct {
	CartID           string        `json:"cartId"`
	BillingAddress   AddressParams `json:"billingAddress"`
	ShippingAddress  AddressParams `json:"shippingAddress"`
	PaymentMethodId  string        `json:"paymentMethodId"`
	ShippingMethodId string        `json:"shippingMethodId"`
	FirstName        string        `json:"firstName"`
	LastName         string        `json:"lastName"`
	Email            string        `json:"email"`
	Note             string        `json:"note"`
	CouponCode       string        `json:"couponCode"`
}

func (GuestCheckoutPlaceOrderParams) GetGraphQLType() string {
	return "GuestCheckoutPlaceOrderParams"
}

type Address struct {
	Street    string   `json:"street"`
	StreetTwo string   `json:"streetTwo"`
	City      string   `json:"city"`
	State     string   `json:"state"`
	Postcode  string   `json:"postcode"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Location  Location `json:"location"`
}

type AddressParams struct {
	Street     string `json:"street"`
	StreetTwo  string `json:"streetTwo"`
	City       string `json:"city"`
	State      string `json:"state"`
	Postcode   string `json:"postcode"`
	LocationId string `json:"locationId"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

func (AddressParams) GetGraphQLType() string {
	return "AddressParams"
}

type Location struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ShippingMethod struct {
	ID                 string `json:"id"`
	DisplayName        string `json:"displayName"`
	DeliveryTimeInDays int    `json:"deliveryTimeInDays"`
	DeliveryCharge     int64  `json:"deliveryCharge"`
	IsFlat             bool   `json:"isFlat"`
}

type PaymentMethod struct {
	ID               string `json:"id"`
	DisplayName      string `json:"displayName"`
	IsDigitalPayment bool   `json:"isDigitalPayment"`
}

type OrderedItem struct {
	PurchasePrice int64   `json:"purchasePrice"`
	Quantity      int     `json:"quantity"`
	Product       Product `json:"product"`
}

type OrderedCart struct {
	CartItems []OrderedItem `json:"cartItems"`
}

type OrderDetail struct {
	ID                   string         `json:"id"`
	Hash                 string         `json:"hash"`
	Subtotal             int64          `json:"subtotal"`
	GrandTotal           int64          `json:"grandTotal"`
	DiscountedAmount     int64          `json:"discountedAmount"`
	Status               string         `json:"status"`
	PaymentStatus        string         `json:"paymentStatus"`
	CreatedAt            string         `json:"createdAt"`
	Cart                 OrderedCart    `json:"cart"`
	CouponCode           CouponCode     `json:"couponCode"`
	Note                 string         `json:"note"`
	ShippingCharge       int64          `json:"shippingCharge"`
	PaymentProcessingFee int64          `json:"paymentProcessingFee"`
	BillingAddress       Address        `json:"billingAddress"`
	ShippingAddress      Address        `json:"shippingAddress"`
	PaymentMethod        PaymentMethod  `json:"paymentMethod"`
	ShippingMethod       ShippingMethod `json:"shippingMethod"`
	Customer             Customer       `json:"customer"`
}

type Customer struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CouponCode struct {
	Code string `json:"code"`
}

type Shop struct {
	Name            string   `json:"name"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	MetaName        string   `json:"metaName"`
	MetaDescription string   `json:"metaDescription"`
	MetaTags        []string `json:"metaTags"`
	Logo            string   `json:"logo"`
	LogoPath        string   `json:"logoPath"`
	Favicon         string   `json:"favicon"`
	FaviconPath     string   `json:"faviconPath"`
	IsOpen          bool     `json:"isOpen"`
	Currency        string   `json:"currency"`
	SupportPhone    string   `json:"supportPhone"`
}

type Cart struct {
	ID                 string          `json:"id"`
	IsShippingRequired bool            `json:"isShippingRequired"`
	CartItems          []QueryCartItem `json:"cartItems"`
}

type PaymentNonce struct {
	PaymentGatewayName   string `json:"PaymentGatewayName" graphql:"PaymentGatewayName"`
	Nonce                string `json:"Nonce" graphql:"Nonce"`
	StripePublishableKey string `json:"StripePublishableKey" graphql:"StripePublishableKey"`
}

type PaymentCallbackOverrides struct {
	SuccessCallback string `json:"SuccessCallback"`
	FailureCallback string `json:"FailureCallback"`
}

func (PaymentCallbackOverrides) GetGraphQLType() string {
	return "PaymentRequestOverrides"
}

type Category struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ProductCount int    `json:"product_count"`
}
