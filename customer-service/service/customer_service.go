package service

import (
	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/customer-service/event"
	"github.com/oliknight1/retail-isa-investment/customer-service/model"
	"github.com/oliknight1/retail-isa-investment/customer-service/repository"
)

type CustomerService interface {
	RegisterCustomer(name string) (model.Customer, error)
}

type customerServiceImpl struct {
	repo      repository.Repository
	publisher event.EventPublisher
}

func New(repo repository.Repository, publisher event.EventPublisher) *customerServiceImpl {
	return &customerServiceImpl{
		repo,
		publisher,
	}
}

func (cs *customerServiceImpl) RegisterCustomer(name string) (model.Customer, error) {
	customer := model.Customer{
		Id:   uuid.New().String(),
		Name: name,
	}

	if err := cs.repo.Create(customer); err != nil {
		return model.Customer{}, err
	}

	if err := cs.publisher.PublishCustomer(customer); err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}
