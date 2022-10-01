package handlers

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (r *Router) handlePayment(ctx *fiber.Ctx) error {
	order, err := r.shopemaaService.OrderDetailsForGuest(ctx.Params("hash"), ctx.Query("email"))
	if err != nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	nonce, err := r.shopemaaService.GeneratePaymentNonce(order.ID, order.Hash, order.Customer.Email, r.cfg.URL)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if nonce.PaymentGatewayName == "SSLCommerz" {
		return ctx.Redirect(nonce.Nonce, http.StatusTemporaryRedirect)
	} else if nonce.PaymentGatewayName == "Stripe" {
		return ctx.Render("payment", map[string]interface{}{
			"order":      order,
			"nonce":      nonce,
			"shop":       r.shopemaaService.GetShop(),
			"page_title": "Payment",
		})
	}

	return ctx.SendStatus(http.StatusBadRequest)
}
