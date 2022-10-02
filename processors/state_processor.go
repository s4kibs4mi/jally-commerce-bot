package processors

import (
	"github.com/s4kibs4mi/jally-commerce-bot/models/api_request"
)

type IStateProcessor interface {
	Init() error
	Process(req *api_request.CustomerRequest) error
	ProcessOrderCreated(cartID, orderHash, email string) error
}
