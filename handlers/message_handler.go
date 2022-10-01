package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/s4kibs4mi/jally-commerce-bot/log"
	"github.com/s4kibs4mi/jally-commerce-bot/models/api_request"
)

func (r *Router) handleMessageReceived(ctx *fiber.Ctx) error {
	req := api_request.FromFiberRequest(ctx)
	if err := r.processor.Process(req); err != nil {
		log.Log().Errorln(err)
	}
	return nil
}
