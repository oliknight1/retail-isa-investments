package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/customer-service/model"
	"github.com/oliknight1/retail-isa-investment/customer-service/repository"
)

func TestCreateCustomerValidation(t *testing.T) {
	db := repository.New()
	validUUID := uuid.New().String()
	tests := []struct {
		name      string
		customer  model.Customer
		expectErr bool
	}{
		{
			name: "valid customer",
			customer: model.Customer{
				Id:   validUUID,
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
				Id:   validUUID,
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
