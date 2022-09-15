package services

import (
	"context"
	"fmt"
	"github.com/hasura/go-graphql-client"
	cfg "github.com/s4kibs4mi/twilfe/config"
	"github.com/s4kibs4mi/twilfe/models"
	"net/http"
)

type IShopemaaService interface {
	GetName() string
	GetCurrency() string
	GetShop() *models.Shop
	ListProducts() ([]models.Product, error)
	AddToCart(productIDs []string) (string, error)
	ConfirmOrder(params *models.PlaceOrderParams) (string, error)
	Pay(orderID string) error
	OrderDetails(orderID string) (*models.OrderDetail, error)
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

func (ss *ShopemaaService) listLocations() ([]models.Location, error) {
	var locationsQuery struct {
		Locations []models.Location `json:"locations"`
	}

	if err := ss.client.Query(context.Background(), &locationsQuery, nil, graphql.OperationName("locations")); err != nil {
		return nil, err
	}

	return locationsQuery.Locations, nil
}

func (ss *ShopemaaService) listShippingMethods() ([]models.ShippingMethod, error) {
	var query struct {
		ShippingMethods []models.ShippingMethod `json:"shippingMethods"`
	}

	if err := ss.client.Query(context.Background(), &query, nil, graphql.OperationName("shippingMethods")); err != nil {
		return nil, err
	}

	return query.ShippingMethods, nil
}

func (ss *ShopemaaService) listPaymentMethods() ([]models.PaymentMethod, error) {
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
	locations, err := ss.listLocations()
	if err != nil {
		return "", err
	}
	paymentMethods, err := ss.listPaymentMethods()
	if err != nil {
		return "", err
	}
	shippingMethods, err := ss.listShippingMethods()
	if err != nil {
		return "", err
	}

	var placeOrderQuery struct {
		OrderGuestCheckout struct {
			Hash string `json:"hash"`
		} `graphql:"orderGuestCheckout(params: $params)"`
	}

	var variables = map[string]interface{}{
		"params": models.GuestCheckoutPlaceOrderParams{
			CartID:    params.CartID,
			FirstName: params.FirstName,
			LastName:  params.LastName,
			Email:     fmt.Sprintf("test@shopemaa.com"),
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

func (ss *ShopemaaService) Pay(orderID string) error {
	return nil
}

func (ss *ShopemaaService) OrderDetails(orderID string) (*models.OrderDetail, error) {
	var orderDetailsQuery struct {
		Order models.OrderDetail `json:"orderByCustomerEmail" graphql:"orderByCustomerEmail(hash: $hash, email: $email)"`
	}

	variables := map[string]interface{}{
		"hash":  graphql.String(orderID),
		"email": graphql.String("test@shopemaa.com"),
	}

	if err := ss.client.Query(context.Background(), &orderDetailsQuery, variables, graphql.OperationName("orderByCustomerEmail")); err != nil {
		return nil, err
	}

	return &orderDetailsQuery.Order, nil
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
