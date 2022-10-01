package services

import (
	"fmt"
	"github.com/s4kibs4mi/jally-commerce-bot/models"
)

type IStateService interface {
	SetState(customerID string, state models.CustomerState) error
	GetState(customerID string) (models.CustomerState, error)
	SetData(customerID, key string, data interface{}) error
	GetData(customer, key string) (interface{}, error)
	SetIdentityByCartID(key, cartID, identity string) error
	GetIdentityByCartID(key, cartID string) (string, error)
}

type StateService struct {
	states map[string]models.CustomerState
	data   map[string]interface{}
}

func (s *StateService) SetState(customerID string, state models.CustomerState) error {
	s.states[customerID] = state
	return nil
}

func (s *StateService) GetState(customerID string) (models.CustomerState, error) {
	return s.states[customerID], nil
}

func (s *StateService) SetData(customerID, key string, data interface{}) error {
	s.data[fmt.Sprintf("%s_%s", customerID, key)] = data
	return nil
}

func (s *StateService) GetData(customerID, key string) (interface{}, error) {
	return s.data[fmt.Sprintf("%s_%s", customerID, key)], nil
}

func (s *StateService) SetIdentityByCartID(key, cartID, identity string) error {
	s.data[fmt.Sprintf("identity_%s_%s", cartID, key)] = identity
	return nil
}

func (s *StateService) GetIdentityByCartID(key, cartID string) (string, error) {
	return s.data[fmt.Sprintf("identity_%s_%s", cartID, key)].(string), nil
}

func NewStateService() IStateService {
	return &StateService{
		states: map[string]models.CustomerState{},
		data:   map[string]interface{}{},
	}
}
