package service_test

import (
	"errors"
	"testing"

	"github.com/oliknight1/retail-isa-investment/customer-service/model"
	"github.com/oliknight1/retail-isa-investment/customer-service/service"
)

type mockRepo struct {
	createFn func(customer model.Customer) error
}

func (m *mockRepo) Create(customer model.Customer) error {
	return m.createFn(customer)
}

func (m *mockRepo) GetById(id string) (*model.Customer, error) {
	return nil, nil
}

type mockPublisher struct {
	publishFn func(customer model.Customer) error
}

func (m *mockPublisher) PublishCustomer(c model.Customer) error {
	return m.publishFn(c)
}

func TestRegisterSuccess(t *testing.T) {
	expectedName := "Oli"
	repo := &mockRepo{
		createFn: func(customer model.Customer) error {
			if customer.Name != expectedName {
				t.Errorf("expected name in register: %s, got: %s", expectedName, customer.Name)
			}
			return nil
		},
	}
	pub := &mockPublisher{
		publishFn: func(customer model.Customer) error {
			if customer.Name != expectedName {
				t.Errorf("expected name in publisher: %s, got %s", expectedName, customer.Name)
			}
			return nil
		},
	}

	svc := service.New(repo, pub)

	customer, err := svc.RegisterCustomer(expectedName)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.Name != expectedName {
		t.Errorf("expected customer name '%s', got %s", expectedName, customer.Name)
	}
	if customer.Id == "" {
		t.Errorf("expected non-empty UUID")
	}
}

func TestRepoFails(t *testing.T) {
	repo := &mockRepo{
		createFn: func(customer model.Customer) error {
			return errors.New("failed to create user")
		},
	}
	pub := &mockPublisher{
		publishFn: func(customer model.Customer) error {
			return nil
		},
	}
	svc := service.New(repo, pub)
	_, err := svc.RegisterCustomer("Oli")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

}
func TestPubFails(t *testing.T) {
	repo := &mockRepo{
		createFn: func(customer model.Customer) error {
			return nil
		},
	}
	pub := &mockPublisher{
		publishFn: func(customer model.Customer) error {
			return errors.New("publish fails")
		},
	}
	svc := service.New(repo, pub)
	_, err := svc.RegisterCustomer("Oli")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

}
