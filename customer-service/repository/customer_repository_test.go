package repository_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/customer-service/model"
	"github.com/oliknight1/retail-isa-investment/customer-service/repository"
)

func TestCreateCustomerValidation(t *testing.T) {
	db := repository.New()
	validId := uuid.New().String()
	tests := []struct {
		name      string
		customer  model.Customer
		expectErr bool
	}{
		{
			name: "valid customer",
			customer: model.Customer{
				Id:   validId,
				Name: "Oli",
			},
			expectErr: false,
		},
		{
			name: "missing ID",
			customer: model.Customer{
				Id:   "",
				Name: "Oli",
			},
			expectErr: true,
		},
		{
			name: "invalid UUID",
			customer: model.Customer{
				Id:   "not-a-uuid",
				Name: "Oli",
			},
			expectErr: true,
		},
		{
			name: "missing name",
			customer: model.Customer{
				Id:   validId,
				Name: "",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Create(tt.customer)
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}

}
func TestGetValidId(t *testing.T) {
	validId := uuid.New().String()
	expectedCustomer := model.Customer{
		Id:   validId,
		Name: "Oli",
	}

	db := &repository.InMemDb{
		Store: map[string]model.Customer{
			validId: expectedCustomer,
		},
	}

	c, err := db.GetById(validId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c == nil {
		t.Fatal("expected customer, got nil")
	}
	if *c != expectedCustomer {
		t.Errorf("got %+v, want %+v", *c, expectedCustomer)
	}
}
func TestGetInvalidUUID(t *testing.T) {
	invalidID := "not-a-uuid"

	db := &repository.InMemDb{
		Store: map[string]model.Customer{},
	}

	c, err := db.GetById(invalidID)
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}
	if c != nil {
		t.Errorf("expected nil customer, got %+v", c)
	}
}
func TestGetByIdNotFound(t *testing.T) {
	missingID := uuid.NewString()

	db := &repository.InMemDb{
		Store: map[string]model.Customer{},
	}

	c, err := db.GetById(missingID)
	if err == nil {
		t.Fatal("expected error for missing customer, got nil")
	}

	expectedErr := fmt.Sprintf("customer with ID %s not found", missingID)
	if err.Error() != expectedErr {
		t.Errorf("unexpected error: got %v, want %v", err.Error(), expectedErr)
	}
	if c != nil {
		t.Errorf("expected nil customer, got %+v", c)
	}
}
