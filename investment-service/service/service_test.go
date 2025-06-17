package service_test

import (
	"testing"

	"github.com/oliknight1/retail-isa-investment/investment-service/model"
	"github.com/oliknight1/retail-isa-investment/investment-service/service"
)

type mockRepo struct {
	createInvestment           func(investment model.Investment) error
	getInvestmentById          func(id string) (*model.Investment, error)
	getInvestmentsByCustomerId func(id string) (*[]model.Investment, error)
}

func (m *mockRepo) CreateInvestment(investment *model.Investment) error {
	return m.createInvestment(*investment)
}
func (m *mockRepo) GetInvestmentById(id string) (*model.Investment, error) {
	return m.getInvestmentById(id)
}

func (m *mockRepo) GetInvestmentsByCustomerId(id string) (*[]model.Investment, error) {
	return m.getInvestmentsByCustomerId(id)
}

type mockPublisher struct {
	publishFn func(subject string, payload any) error
	close     func()
}

func (m *mockPublisher) Publish(subject string, payload any) error {
	return m.publishFn(subject, payload)
}
func (m *mockPublisher) Close() {
	m.close()
}
func TestCreateInvestment(t *testing.T) {
	tests := []struct {
		name       string
		customerId string
		fundId     string
		amount     float64
		expectErr  bool
		repoFn     func(model.Investment) error
		publishFn  func(string, any) error
	}{
		{
			name:       "success",
			customerId: "cust-1",
			fundId:     "fund-1",
			amount:     100,
			expectErr:  false,
			repoFn: func(inv model.Investment) error {
				return nil
			},
			publishFn: func(subject string, payload any) error {
				return nil
			},
		},
		{
			name:       "error - empty customerId",
			customerId: "",
			fundId:     "fund-1",
			amount:     100,
			expectErr:  true,
		},
		{
			name:       "error - empty fundId",
			customerId: "cust-1",
			fundId:     "",
			amount:     100,
			expectErr:  true,
		},
		{
			name:       "error - zero amount",
			customerId: "cust-1",
			fundId:     "fund-1",
			amount:     0,
			expectErr:  true,
		},
		{
			name:       "error - negative amount",
			customerId: "cust-1",
			fundId:     "fund-1",
			amount:     -50,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepo{
				createInvestment: func(inv model.Investment) error {
					if tt.repoFn != nil {
						return tt.repoFn(inv)
					}
					return nil
				},
			}
			mockPub := &mockPublisher{
				publishFn: func(subject string, payload any) error {
					if tt.publishFn != nil {
						return tt.publishFn(subject, payload)
					}
					return nil
				},
				close: func() {},
			}

			svc := service.New(mockRepo, mockPub)
			inv, err := svc.CreateInvestment(tt.customerId, tt.fundId, tt.amount)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				if inv != nil {
					t.Errorf("expected nil investment, got: %+v", inv)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if inv == nil {
				t.Errorf("expected investment, got nil")
				return
			}
			if inv.CustomerId != tt.customerId || inv.FundId != tt.fundId || inv.Amount != tt.amount {
				t.Errorf("unexpected investment data: %+v", inv)
			}
		})
	}
}
