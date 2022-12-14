package services

import (
	"github.com/s4kibs4mi/jally-commerce-bot/config"
	twilio2 "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type ITwilioService interface {
	Send(to, from, payload string) error
	Api() *twilio2.RestClient
}

type TwilioService struct {
	client *twilio2.RestClient
}

func (ts *TwilioService) Send(to, from, payload string) error {
	params := &openapi.CreateMessageParams{
		To:   &to,
		From: &from,
		Body: &payload,
	}
	_, err := ts.client.Api.CreateMessage(params)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TwilioService) Api() *twilio2.RestClient {
	return ts.client
}

func NewTwilioService(cfg *config.Application) ITwilioService {
	c := twilio2.NewRestClientWithParams(twilio2.ClientParams{
		Username: cfg.TwilioUsername,
		Password: cfg.TwilioPassword,
	})

	return &TwilioService{
		client: c,
	}
}
