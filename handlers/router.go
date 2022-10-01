package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/s4kibs4mi/jally-commerce-bot/config"
	"github.com/s4kibs4mi/jally-commerce-bot/processors"
	"github.com/s4kibs4mi/jally-commerce-bot/services"
	"time"
)

type IRouter interface {
	Start(addr string) error
	Stop() error
}

type Router struct {
	cfg             *config.Application
	app             *fiber.App
	processor       processors.IStateProcessor
	shopemaaService services.IShopemaaService
}

func (r *Router) Start(addr string) error {
	engine := html.New("./web", ".html")
	engine.AddFunc("currentYear", func() int { return time.Now().Year() })
	engine.AddFunc("formatAmount", func(amount int64) string {
		if amount == 0 {
			return "Free"
		}
		return fmt.Sprintf("%.2f %s", float64(amount)/float64(100), r.shopemaaService.GetCurrency())
	})
	engine.AddFunc("formatAmountR", func(amount int64) string {
		if amount == 0 {
			return fmt.Sprintf("%.2f %s", float64(0), r.shopemaaService.GetCurrency())
		}
		return fmt.Sprintf("%.2f %s", float64(amount)/float64(100), r.shopemaaService.GetCurrency())
	})
	engine.AddFunc("getIndex", func(arr []string, index int) string { return arr[index] })

	r.app = fiber.New(fiber.Config{
		Views: engine,
	})

	r.app.Use(logger.New())
	r.app.Use(recover.New())
	r.app.Static("/static", "./web/assets")
	r.app.Get("/", r.handleIndex)
	r.app.Post("/messages/", r.handleMessageReceived)
	r.app.Get("/orders/:hash/", r.handleOrderDetails)
	r.app.Post("/orders/:hash/", r.handleOrderDetails)
	r.app.Get("/orders/:hash/payment/", r.handlePayment)
	r.app.Get("/checkout/:cartId/", r.handleCheckout)
	r.app.Post("/checkout/:cartId/", r.handleCheckoutSubmit)
	return r.app.Listen(addr)
}

func (r *Router) Stop() error {
	return r.app.Shutdown()
}

func NewRouter(cfg *config.Application, processor processors.IStateProcessor, shopemaaService services.IShopemaaService) IRouter {
	return &Router{
		cfg:             cfg,
		processor:       processor,
		shopemaaService: shopemaaService,
	}
}
