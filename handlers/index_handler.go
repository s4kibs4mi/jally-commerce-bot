package handlers

import "github.com/gofiber/fiber/v2"

func (r *Router) handleIndex(ctx *fiber.Ctx) error {
	return ctx.Send([]byte("Welcome to Twilfe"))
}
