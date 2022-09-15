package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func (r *Router) handleOrderDetails(ctx *fiber.Ctx) error {
	order, err := r.shopemaaService.OrderDetails(ctx.Params("hash"))
	if err != nil {
		return ctx.SendStatus(404)
	}

	var items HTML

	for _, item := range order.Cart.CartItems {
		items += Tr_(
			Td(Attr(a.Class_("text-center")), Text(item.Product.Name)),
			Td(Attr(a.Class_("text-center")), Text(fmt.Sprintf("%.2f %s", float64(item.PurchasePrice)/float64(100), r.shopemaaService.GetCurrency()))),
		)
	}

	page := Html5_(
		Head_(
			Link(Attr(
				a.Href_("https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/css/bootstrap.min.css"),
				a.Rel("stylesheet"),
			)),
			Script(
				Attr(
					a.Src("https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/js/bootstrap.min.js"),
					a.Type_("javascript"),
				),
				JavaScript_(),
			),
			Title_(HTML(fmt.Sprintf("Order#%s - %s", order.Hash, r.shopemaaService.GetName()))),
			Link(
				Attr(
					a.Rel_("icon"),
					a.Type_("image/x-icon"),
					a.Href_(r.shopemaaService.GetShop().Logo),
				),
			),
			Meta(
				Attr(a.Name_("og:title"), a.Value_(fmt.Sprintf("Order#%s - %s", order.Hash, r.shopemaaService.GetShop().Name))),
			),
			Meta(
				Attr(a.Name_("og:type"), a.Value_("website")),
			),
			Meta(
				Attr(a.Name_("og:image"), a.Value_(r.shopemaaService.GetShop().Logo)),
			),
		),
		Body_(
			H1(Attr(a.Class_("text-center mt-5")), Text("Welcome to "+r.shopemaaService.GetName())),
			Div(
				Attr(a.Class_("container col-12 col-md-12 table-responsive-lg")),
				Table(
					Attr(
						a.Class_("table mt-5"),
					),
					Tr_(
						Td(Attr(a.Class_("text-left text-primary")), Text("Order Id")),
						Td(Attr(a.Class_("text-end")), Text(order.Hash)),
					),
					Tr_(
						Td(Attr(a.Class_("text-left text-primary")), Text("Order Date")),
						Td(Attr(a.Class_("text-end")), Text(order.CreatedAt)),
					),
					Tr_(
						Td(Attr(a.Class_("text-left text-primary")), Text("Order Status")),
						Td(Attr(a.Class_("text-end")), Text(order.Status)),
					),
					Tr_(
						Td(Attr(a.Class_("text-left text-primary")), Text("Payment Status")),
						Td(Attr(a.Class_("text-end")), Text(order.PaymentStatus)),
					),
					Tr_(
						Td(Attr(a.Class_("text-left text-primary")), Text("Subtotal")),
						Td(Attr(a.Class_("text-end")), Text(fmt.Sprintf("%.2f %s", float64(order.Subtotal)/float64(100), r.shopemaaService.GetCurrency()))),
					),
					Tr_(
						Td(Attr(a.Class_("text-left text-primary")), Text("Grand Total")),
						Td(Attr(a.Class_("text-end")), Text(fmt.Sprintf("%.2f %s", float64(order.GrandTotal)/float64(100), r.shopemaaService.GetCurrency()))),
					),
					Tr_(
						Td(Attr(a.Class_("text-left text-primary mt-5")), Text("Items")),
						Td(Attr(a.Class_("text-end text-primary mt-5")), Text("")),
					),
					items,
				),
			),
		))

	ctx.Response().Header.Set("Content-Type", "text/html")
	return ctx.SendString(string(page))
}
