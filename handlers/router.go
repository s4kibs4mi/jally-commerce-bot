package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/s4kibs4mi/jally-commerce-bot/config"
	"github.com/s4kibs4mi/jally-commerce-bot/log"
	"github.com/s4kibs4mi/jally-commerce-bot/models/api_request"
	"github.com/s4kibs4mi/jally-commerce-bot/processors"
	"github.com/s4kibs4mi/jally-commerce-bot/services"
	"github.com/s4kibs4mi/jally-commerce-bot/services/messenger"
	"net/http"
	"time"
)

type IRouter interface {
	Start(addr string) error
	Stop() error
}

type Router struct {
	cfg               *config.Application
	app               *fiber.App
	processor         processors.IStateProcessor
	shopemaaService   services.IShopemaaService
	facebookMessenger *messenger.Messenger
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
	r.app.Get("/orders/:hash/", r.handleOrderDetails)
	r.app.Post("/orders/:hash/", r.handleOrderDetails)
	r.app.Get("/orders/:hash/payment/", r.handlePayment)
	r.app.Get("/checkout/:cartId/", r.handleCheckout)
	r.app.Post("/checkout/:cartId/", r.handleCheckoutSubmit)

	r.app.Post("/messages/", r.handleMessageReceived) // Twilio

	if r.facebookMessenger != nil {
		facebookGroup := r.app.Group("/facebook")
		facebookGroup.Use(func(ctx *fiber.Ctx) error {
			// Authenticate
			return ctx.Next()
		})
		facebookGroup.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(ctx.Query("hub.challenge"))
		})
		facebookGroup.Post("/", r.facebookMessenger.Handler)
		r.facebookMessenger.MessageReceived = func(event messenger.Event, opts messenger.MessageOpts, message messenger.ReceivedMessage) {
			if err := r.processor.Process(&api_request.CustomerRequest{
				Event:     event,
				Opts:      opts,
				Message:   message,
				IsMessage: true,
			}); err != nil {
				log.Log().Errorln(err)
			}
		}
		r.facebookMessenger.Postback = func(event messenger.Event, opts messenger.MessageOpts, postback messenger.Postback) {
			if err := r.processor.Process(&api_request.CustomerRequest{
				Event:     event,
				Opts:      opts,
				Postback:  postback,
				IsMessage: false,
			}); err != nil {
				log.Log().Errorln(err)
			}
		}
	}

	return r.app.Listen(addr)
}

func (r *Router) Stop() error {
	return r.app.Shutdown()
}

func NewRouter(cfg *config.Application, processor processors.IStateProcessor,
	shopemaaService services.IShopemaaService, facebookMessenger *messenger.Messenger) IRouter {
	return &Router{
		cfg:               cfg,
		processor:         processor,
		shopemaaService:   shopemaaService,
		facebookMessenger: facebookMessenger,
	}
}
