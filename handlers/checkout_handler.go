package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/s4kibs4mi/twilfe/log"
	"github.com/s4kibs4mi/twilfe/models"
	"sync"
)

func (r *Router) handleCheckout(ctx *fiber.Ctx) error {
	shop := r.shopemaaService.GetShop()
	var shippingMethods []models.ShippingMethod
	var paymentMethods []models.PaymentMethod
	var locations []models.Location
	var cart *models.Cart
	total := int64(0)
	discount := int64(0)

	wg := sync.WaitGroup{}
	wg.Add(5)

	go func() {
		sm, err := r.shopemaaService.ListShippingMethods()
		if err == nil {
			shippingMethods = sm
		}
		wg.Done()
	}()
	go func() {
		pm, err := r.shopemaaService.ListPaymentMethods()
		if err == nil {
			paymentMethods = pm
		}
		wg.Done()
	}()
	go func() {
		l, err := r.shopemaaService.ListLocations()
		if err == nil {
			locations = l
		}
		wg.Done()
	}()
	go func() {
		c, err := r.shopemaaService.GetCart(ctx.Params("cartId"))
		if err == nil {
			cart = c

			for _, item := range cart.CartItems {
				total += int64(item.Quantity) * item.PurchasePrice
			}
		}
		wg.Done()
	}()
	go func() {
		c, err := r.shopemaaService.CheckDiscount(ctx.Params("cartId"), ctx.Query("coupon"), nil)
		if err == nil {
			discount = c
		}
		wg.Done()
	}()

	wg.Wait()

	return ctx.Render("checkout", map[string]interface{}{
		"shop":            shop,
		"shippingMethods": shippingMethods,
		"paymentMethods":  paymentMethods,
		"locations":       locations,
		"cart":            cart,
		"subTotal":        total,
		"discount":        discount,
		"grandTotal":      total - discount,
		"coupon":          ctx.Query("coupon"),
		"page_title":      "Checkout",
	})
}

func (r *Router) handleCheckoutSubmit(ctx *fiber.Ctx) error {
	pld := struct {
		FirstName      string `json:"firstName"`
		LastName       string `json:"lastName"`
		Phone          string `json:"phone"`
		Email          string `json:"email"`
		ShippingMethod string `json:"shippingMethod"`
		PaymentMethod  string `json:"paymentMethod"`
		Address        string `json:"address"`
		Address2       string `json:"address2"`
		PostalCode     string `json:"postalCode"`
		City           string `json:"city"`
		State          string `json:"state"`
		Country        string `json:"country"`
		Note           string `json:"note"`
	}{}
	if err := ctx.BodyParser(&pld); err != nil {
		return r.handleCheckout(ctx)
	}

	orderHash, err := r.shopemaaService.PlaceOrder(&models.GuestCheckoutPlaceOrderParams{
		FirstName:        pld.FirstName,
		LastName:         pld.LastName,
		Email:            pld.Email,
		CartID:           ctx.Params("cartId"),
		PaymentMethodId:  pld.PaymentMethod,
		ShippingMethodId: pld.ShippingMethod,
		ShippingAddress: models.AddressParams{
			Street:     pld.Address,
			StreetTwo:  pld.Address2,
			City:       pld.City,
			State:      pld.State,
			Postcode:   pld.PostalCode,
			LocationId: pld.Country,
			Email:      pld.Email,
			Phone:      pld.Phone,
		},
		BillingAddress: models.AddressParams{
			Street:     pld.Address,
			StreetTwo:  pld.Address2,
			City:       pld.City,
			State:      pld.State,
			Postcode:   pld.PostalCode,
			LocationId: pld.Country,
			Email:      pld.Email,
			Phone:      pld.Phone,
		},
	})
	if err != nil {
		return r.handleCheckout(ctx)
	}

	cartID := ctx.Params("cartId")
	if err := r.processor.ProcessOrderCreated(cartID, orderHash, pld.Email); err != nil {
		log.Log().Errorln(err)
	}

	return ctx.Redirect(fmt.Sprintf("/orders/%s/?email=%s", orderHash, pld.Email))
}
