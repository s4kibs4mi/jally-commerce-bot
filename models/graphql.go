package models

type Search struct {
	Query   string   `json:"query"`
	Filters []string `json:"filters"`
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
}

func (GuestCheckoutPlaceOrderParams) GetGraphQLType() string {
	return "GuestCheckoutPlaceOrderParams"
}

type AddressParams struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	Postcode   string `json:"postcode"`
	LocationId string `json:"locationId"`
}

func (AddressParams) GetGraphQLType() string {
	return "AddressParams"
}

type Location struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ShippingMethod struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type PaymentMethod struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type OrderedItem struct {
	PurchasePrice int64   `json:"purchasePrice"`
	Product       Product `json:"product"`
}

type OrderedCart struct {
	CartItems []OrderedItem `json:"cartItems"`
}

type OrderDetail struct {
	ID               string      `json:"id"`
	Hash             string      `json:"hash"`
	Subtotal         int64       `json:"subtotal"`
	GrandTotal       int64       `json:"grandTotal"`
	DiscountedAmount int64       `json:"discountedAmount"`
	Status           string      `json:"status"`
	PaymentStatus    string      `json:"paymentStatus"`
	CreatedAt        string      `json:"createdAt"`
	Cart             OrderedCart `json:"cart"`
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
}
