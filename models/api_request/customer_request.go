package api_request

import (
	"github.com/gofiber/fiber/v2"
	"github.com/s4kibs4mi/jally-commerce-bot/services/messenger"
)

type CustomerRequest struct {
	AccountSid  string
	Body        string
	From        string
	To          string
	MessageSid  string
	ProfileName string

	Event     messenger.Event
	Opts      messenger.MessageOpts
	Message   messenger.ReceivedMessage
	Postback  messenger.Postback
	IsMessage bool
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
