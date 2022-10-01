package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hasura/go-graphql-client"
	cfg "github.com/s4kibs4mi/twilfe/config"
	"github.com/s4kibs4mi/twilfe/models"
	"net/http"
	"net/url"
)

type IShopemaaService interface {
	GetName() string
	GetCurrency() string
	GetShop() *models.Shop
	ListProducts() ([]models.Product, error)
	AddToCart(productIDs []string) (string, error)
	GetCart(cartID string) (*models.Cart, error)
	ConfirmOrder(params *models.PlaceOrderParams) (string, error)
	PlaceOrder(params *models.GuestCheckoutPlaceOrderParams) (string, error)
	Pay(orderID string) error
	OrderDetails(orderID string) (*models.OrderDetail, error)
	OrderDetailsForGuest(orderHash, email string) (*models.OrderDetail, error)
	ListShippingMethods() ([]models.ShippingMethod, error)
	ListPaymentMethods() ([]models.PaymentMethod, error)
	ListLocations() ([]models.Location, error)
	CheckDiscount(cartID, couponCode string, shippingMethodID *string) (int64, error)
}

type ShopemaaService struct {
	client *graphql.Client
	shop   *models.Shop
}

func (ss *ShopemaaService) GetName() string {
	return ss.shop.Name
}

func (ss *ShopemaaService) GetCurrency() string {
	return ss.shop.Currency
}

func (ss *ShopemaaService) GetShop() *models.Shop {
	b, _ := url.QueryUnescape(ss.shop.Description)
	ss.shop.Description = b
	return ss.shop
}

func (ss *ShopemaaService) ListProducts() ([]models.Product, error) {
	var productsQuery struct {
		Products []models.Product `json:"products" graphql:"products(search: $search, sort: $sort, pagination: $pagination)"`
	}

	variables := map[string]interface{}{
		"search": models.Search{
			Query:   "",
			Filters: []string{},
		},
		"sort": models.Sort{
			By:        "CreatedAt",
			Direction: "Desc",
		},
		"pagination": models.Pagination{
			PerPage: 25,
			Page:    1,
		},
	}

	if err := ss.client.Query(context.Background(), &productsQuery, variables, graphql.OperationName("products")); err != nil {
		return nil, err
	}

	return productsQuery.Products, nil
}

func (ss *ShopemaaService) GetCart(cartID string) (*models.Cart, error) {
	var cartQuery struct {
		Cart models.Cart `graphql:"cart(cartId: $cartId)"`
	}

	var variables = map[string]interface{}{
		"cartId": graphql.String(cartID),
	}

	if err := ss.client.Query(context.Background(), &cartQuery, variables, graphql.OperationName("cart")); err != nil {
		return nil, err
	}

	return &cartQuery.Cart, nil
}

func (ss *ShopemaaService) createEmptyCart() (string, error) {
	var cartMutation struct {
		NewCart struct {
			ID string `json:"id"`
		} `graphql:"newCart(params: $params)"`
	}

	var variables = map[string]interface{}{
		"params": models.NewCartParams{
			CartItems: []models.CartItem{},
		},
	}

	if err := ss.client.Mutate(context.Background(), &cartMutation, variables, graphql.OperationName("newCart")); err != nil {
		return "", err
	}

	return cartMutation.NewCart.ID, nil
}

func (ss *ShopemaaService) ListLocations() ([]models.Location, error) {
	var locationsQuery struct {
		Locations []models.Location `json:"locations"`
	}

	if err := ss.client.Query(context.Background(), &locationsQuery, nil, graphql.OperationName("locations")); err != nil {
		return nil, err
	}

	return locationsQuery.Locations, nil
}

func (ss *ShopemaaService) ListShippingMethods() ([]models.ShippingMethod, error) {
	var query struct {
		ShippingMethods []models.ShippingMethod `json:"shippingMethods"`
	}

	if err := ss.client.Query(context.Background(), &query, nil, graphql.OperationName("shippingMethods")); err != nil {
		return nil, err
	}

	return query.ShippingMethods, nil
}

func (ss *ShopemaaService) ListPaymentMethods() ([]models.PaymentMethod, error) {
	var query struct {
		PaymentMethods []models.PaymentMethod `json:"paymentMethods"`
	}

	if err := ss.client.Query(context.Background(), &query, nil, graphql.OperationName("paymentMethods")); err != nil {
		return nil, err
	}

	return query.PaymentMethods, nil
}

func (ss *ShopemaaService) AddToCart(productIDs []string) (string, error) {
	cartID, err := ss.createEmptyCart()
	if err != nil {
		return "", err
	}

	var cartQuery struct {
		UpdateCart struct {
			ID string `json:"id"`
		} `graphql:"updateCart(id: $id, params: $params)"`
	}

	var cartItems []models.CartItem
	for _, id := range productIDs {
		cartItems = append(cartItems, models.CartItem{
			ProductID: id,
			Quantity:  1,
		})
	}

	var variables = map[string]interface{}{
		"params": models.UpdateCartParams{
			CartItems: cartItems,
		},
		"id": graphql.String(cartID),
	}

	if err := ss.client.Mutate(context.Background(), &cartQuery, variables, graphql.OperationName("updateCart")); err != nil {
		return cartID, err
	}

	return cartID, nil
}

func (ss *ShopemaaService) ConfirmOrder(params *models.PlaceOrderParams) (string, error) {
	locations, err := ss.ListLocations()
	if err != nil {
		return "", err
	}
	paymentMethods, err := ss.ListShippingMethods()
	if err != nil {
		return "", err
	}
	shippingMethods, err := ss.ListPaymentMethods()
	if err != nil {
		return "", err
	}

	var placeOrderQuery struct {
		OrderGuestCheckout struct {
			Hash string `json:"hash"`
		} `graphql:"orderGuestCheckout(params: $params)"`
	}

	if params.Email == "" {
		params.Email = fmt.Sprintf("test@shopemaa.com")
	}

	var variables = map[string]interface{}{
		"params": models.GuestCheckoutPlaceOrderParams{
			CartID:    params.CartID,
			FirstName: params.FirstName,
			LastName:  params.LastName,
			Email:     params.Email,
			BillingAddress: models.AddressParams{
				Street:     "test",
				City:       "test",
				Postcode:   "test",
				LocationId: locations[0].ID,
			},
			ShippingAddress: models.AddressParams{
				Street:     "test",
				City:       "test",
				Postcode:   "test",
				LocationId: locations[0].ID,
			},
			ShippingMethodId: shippingMethods[0].ID,
			PaymentMethodId:  paymentMethods[0].ID,
		},
	}

	if err := ss.client.Mutate(context.Background(), &placeOrderQuery, variables, graphql.OperationName("orderGuestCheckout")); err != nil {
		return "", err
	}

	return placeOrderQuery.OrderGuestCheckout.Hash, nil
}

func (ss *ShopemaaService) PlaceOrder(params *models.GuestCheckoutPlaceOrderParams) (string, error) {
	var placeOrderQuery struct {
		OrderGuestCheckout struct {
			Hash string `json:"hash"`
		} `graphql:"orderGuestCheckout(params: $params)"`
	}

	var variables = map[string]interface{}{
		"params": params,
	}

	if err := ss.client.Mutate(context.Background(), &placeOrderQuery, variables, graphql.OperationName("orderGuestCheckout")); err != nil {
		b, _ := json.Marshal(err)
		fmt.Println(string(b))
		return "", err
	}

	return placeOrderQuery.OrderGuestCheckout.Hash, nil
}

func (ss *ShopemaaService) Pay(orderID string) error {
	return nil
}

func (ss *ShopemaaService) OrderDetails(orderID string) (*models.OrderDetail, error) {
	return ss.OrderDetailsForGuest(orderID, "test@shopemaa.com")
}

func (ss *ShopemaaService) OrderDetailsForGuest(orderHash, email string) (*models.OrderDetail, error) {
	var orderDetailsQuery struct {
		Order models.OrderDetail `json:"orderByCustomerEmail" graphql:"orderByCustomerEmail(hash: $hash, email: $email)"`
	}

	variables := map[string]interface{}{
		"hash":  graphql.String(orderHash),
		"email": graphql.String(email),
	}

	if err := ss.client.Query(context.Background(), &orderDetailsQuery, variables, graphql.OperationName("orderByCustomerEmail")); err != nil {
		return nil, err
	}

	return &orderDetailsQuery.Order, nil
}

func (ss *ShopemaaService) CheckDiscount(cartID, couponCode string, shippingMethodID *string) (int64, error) {
	var checkDiscountMutation struct {
		Discount int64 `graphql:"checkDiscountForGuests(couponCode: $couponCode, cartId: $cartId)"`
	}

	var variables = map[string]interface{}{
		"cartId":     graphql.String(cartID),
		"couponCode": graphql.String(couponCode),
	}
	//if shippingMethodID != nil {
	//	variables["shippingMethodID"] = graphql.String(*shippingMethodID)
	//} else {
	//	variables["shippingMethodID"] = graphql.
	//}

	if err := ss.client.Query(context.Background(), &checkDiscountMutation, variables, graphql.OperationName("checkDiscountForGuests")); err != nil {
		return 0, err
	}

	return checkDiscountMutation.Discount, nil
}

func NewShopemaaService(cfg *cfg.Application) (IShopemaaService, error) {
	c := graphql.
		NewClient("https://api.shopemaa.com/query", &http.Client{}).
		WithDebug(true).
		WithRequestModifier(func(request *http.Request) {
			request.Header.Set("store-key", cfg.ShopemaaKey)
			request.Header.Set("store-secret", cfg.ShopemaaSecret)
		})

	var shopQuery struct {
		StoreBySecret models.Shop `json:"storeBySecret"`
	}

	err := c.Query(context.Background(), &shopQuery, nil)
	if err != nil {
		return nil, err
	}

	ss := &ShopemaaService{
		client: c,
		shop:   &shopQuery.StoreBySecret,
	}
	return ss, nil
}
