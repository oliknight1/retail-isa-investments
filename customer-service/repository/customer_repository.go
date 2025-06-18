package repository

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/customer-service/model"
)

type Repository interface {
	Create(customer model.Customer) error
	GetById(id string) (*model.Customer, error)
}

type InMemDb struct {
	Store map[string]model.Customer
}

func New() *InMemDb {
	return &InMemDb{
		Store: make(map[string]model.Customer),
	}
}

func (db *InMemDb) Create(customer model.Customer) error {
	if customer.Id == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}
	if customer.Name == "" {
		return fmt.Errorf("customer name cannot be empty")
	}

	if err := uuid.Validate(customer.Id); err != nil {
		return fmt.Errorf("invalid customer ID: %w", err)
	}

	db.Store[customer.Id] = customer
	return nil
}

func (db *InMemDb) GetById(id string) (*model.Customer, error) {
	if err := uuid.Validate(id); err != nil {
		log.Printf("invalid UUID provided: %s, error: %v", id, err)
		return nil, err
	}
	c, ok := db.Store[id]
	if !ok {
		err := fmt.Errorf("customer with ID %s not found", id)
		log.Println(err)
		return nil, err
	}
	return &c, nil
}
