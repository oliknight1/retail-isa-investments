package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/oliknight1/retail-isa-investment/investment-service/internal"
	"github.com/oliknight1/retail-isa-investment/investment-service/logger"
	"github.com/oliknight1/retail-isa-investment/investment-service/model"
	"github.com/oliknight1/retail-isa-investment/investment-service/service"
)

type mockRepo struct {
	createInvestment           func(investment model.Investment) error
	getInvestmentById          func(id string) (*model.Investment, error)
	getInvestmentsByCustomerId func(id string) (*[]model.Investment, error)
}

func (m *mockRepo) CreateInvestment(investment model.Investment) error {
	return m.createInvestment(investment)
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

func TestCreateInvestmentSuccess(t *testing.T) {
	mockRepo := &mockRepo{
		createInvestment: func(inv model.Investment) error {
			return nil
		},
	}
	mockPub := &mockPublisher{
		publishFn: func(subject string, payload any) error {
			return nil
		},
		close: func() {},
	}

	logger := logger.NewMockLogger()
	svc := service.New(mockRepo, mockPub, logger)

	customerId := "cust-1"
	fundId := "fund-1"
	amount := 100.0

	investment, err := svc.CreateInvestment(customerId, fundId, amount)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if investment == nil {
		t.Fatal("expected investment, got nil")
	}

	expected := &model.Investment{
		Id:         investment.Id,
		CustomerId: customerId,
		FundId:     fundId,
		Amount:     amount,
		Status:     "pending",
		CreatedAt:  investment.CreatedAt,
	}

	if diff := cmp.Diff(expected, investment,
		cmp.AllowUnexported(model.Investment{}),
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.String() == "CompletedAt" || p.String() == "FailureReason"
		}, cmp.Ignore()),
	); diff != "" {
		t.Errorf("unexpected investment (-want +got):\n%s", diff)
	}
}

func TestCreateInvestmentFailures(t *testing.T) {
	tests := []struct {
		name        string
		customerId  string
		fundId      string
		amount      float64
		expectedErr string
	}{
		{
			name:        "missing customerId",
			customerId:  "",
			fundId:      "fund-1",
			amount:      100,
			expectedErr: internal.ErrMissingCustomerId.Error(),
		},
		{
			name:        "missing fundId",
			customerId:  "cust-1",
			fundId:      "",
			amount:      100,
			expectedErr: internal.ErrMissingFundId.Error(),
		},
		{
			name:        "amount is zero",
			customerId:  "cust-1",
			fundId:      "fund-1",
			amount:      0,
			expectedErr: internal.ErrZeroTransactionAmount.Error(),
		},
		{
			name:        "amount is negative",
			customerId:  "cust-1",
			fundId:      "fund-1",
			amount:      -50,
			expectedErr: internal.ErrZeroTransactionAmount.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepo{
				createInvestment: func(inv model.Investment) error {
					t.Fatal("should not call CreateInvestment on validation failure")
					return nil
				},
			}
			mockPub := &mockPublisher{
				publishFn: func(subject string, payload any) error {
					t.Fatal("should not publish on validation failure")
					return nil
				},
				close: func() {},
			}

			logger := logger.NewMockLogger()
			svc := service.New(mockRepo, mockPub, logger)

			investment, err := svc.CreateInvestment(tt.customerId, tt.fundId, tt.amount)

			if err == nil {
				t.Fatalf("expected error '%s', got nil", tt.expectedErr)
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("expected error: '%s', got: '%s'", tt.expectedErr, err.Error())
			}
			if investment != nil {
				t.Errorf("expected nil investment, got: %+v", investment)
			}
		})
	}
}

func TestGetInvestmentByIdSuccess(t *testing.T) {
	expected := &model.Investment{
		Id:         "inv-1",
		CustomerId: "cust-1",
		FundId:     "fund-1",
		Amount:     200.0,
		Status:     "completed",
		CreatedAt:  time.Now(),
	}

	mockRepo := &mockRepo{
		getInvestmentById: func(id string) (*model.Investment, error) {
			return expected, nil
		},
	}
	logger := logger.NewMockLogger()
	svc := service.New(mockRepo, nil, logger)

	actual, err := svc.GetInvestmentById("inv-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(expected, actual, cmp.AllowUnexported(model.Investment{})); diff != "" {
		t.Errorf("unexpected result (-want +got):\n%s", diff)
	}
}

func TestGetInvestmentByIdRepoFails(t *testing.T) {
	mockRepo := &mockRepo{
		getInvestmentById: func(id string) (*model.Investment, error) {
			return nil, errors.New("not found")
		},
	}
	logger := logger.NewMockLogger()
	svc := service.New(mockRepo, nil, logger)

	_, err := svc.GetInvestmentById("missing-id")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "not found" {
		t.Errorf("expected error 'not found', got %v", err)
	}
}

func TestGetInvestmentByCustomerIdSuccess(t *testing.T) {
	expected := &[]model.Investment{
		{
			Id:         "inv-1",
			CustomerId: "cust-1",
			FundId:     "fund-1",
			Amount:     200.0,
			Status:     "completed",
			CreatedAt:  time.Now(),
		},
		{
			Id:         "inv-2",
			CustomerId: "cust-1",
			FundId:     "fund-2",
			Amount:     100.0,
			Status:     "completed",
			CreatedAt:  time.Now(),
		},
		{
			Id:         "inv-3",
			CustomerId: "cust-2",
			FundId:     "fund-3",
			Amount:     300.0,
			Status:     "completed",
			CreatedAt:  time.Now(),
		},
	}

	mockRepo := &mockRepo{
		getInvestmentsByCustomerId: func(id string) (*[]model.Investment, error) {
			return expected, nil
		},
	}
	logger := logger.NewMockLogger()
	svc := service.New(mockRepo, nil, logger)

	actual, err := svc.GetInvestmentsByCustomerId("cust-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*actual) != len(*expected) {
		t.Fatalf("expect length of: %d, recieved: %d", len(*expected), len(*actual))
	}
	if diff := cmp.Diff(expected, actual, cmp.AllowUnexported(model.Investment{})); diff != "" {
		t.Errorf("unexpected result (-want +got):\n%s", diff)
	}
}
