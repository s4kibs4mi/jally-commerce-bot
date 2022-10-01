package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) handleOrderDetails(ctx *fiber.Ctx) error {
	order, err := r.shopemaaService.OrderDetailsForGuest(ctx.Params("hash"), ctx.Query("email"))
	if err != nil {
		return ctx.SendStatus(404)
	}
	return ctx.Render("order_details", map[string]interface{}{
		"order":            order,
		"shop":             r.shopemaaService.GetShop(),
		"page_title":       "Order Details",
		"show_payment_btn": order.PaymentStatus != "Paid" && order.PaymentMethod.IsDigitalPayment,
	})
}
