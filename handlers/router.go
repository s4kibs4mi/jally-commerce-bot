package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/s4kibs4mi/twilfe/processors"
	"github.com/s4kibs4mi/twilfe/services"
)

type IRouter interface {
	Start(addr string) error
	Stop() error
}

type Router struct {
	app             *fiber.App
	processor       processors.IStateProcessor
	shopemaaService services.IShopemaaService
}

func (r *Router) Start(addr string) error {
	r.app.Use(logger.New())
	r.app.Get("/", r.handleIndex)
	r.app.Post("/messages/", r.handleMessageReceived)
	r.app.Get("/orders/:hash/", r.handleOrderDetails)
	return r.app.Listen(addr)
}

func (r *Router) Stop() error {
	return r.app.Shutdown()
}

func NewRouter(processor processors.IStateProcessor, shopemaaService services.IShopemaaService) IRouter {
	return &Router{
		app:             fiber.New(),
		processor:       processor,
		shopemaaService: shopemaaService,
	}
}
