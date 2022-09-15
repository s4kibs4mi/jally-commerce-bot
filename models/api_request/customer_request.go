package api_request

import "github.com/gofiber/fiber/v2"

type CustomerRequest struct {
	AccountSid  string
	Body        string
	From        string
	To          string
	MessageSid  string
	ProfileName string
}

func FromFiberRequest(ctx *fiber.Ctx) *CustomerRequest {
	return &CustomerRequest{
		AccountSid:  ctx.FormValue("AccountSid", ""),
		Body:        ctx.FormValue("Body", ""),
		From:        ctx.FormValue("From", ""),
		MessageSid:  ctx.FormValue("MessageSid", ""),
		ProfileName: ctx.FormValue("ProfileName", ""),
		To:          ctx.FormValue("To", ""),
	}
}
