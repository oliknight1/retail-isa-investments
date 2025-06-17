package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/oliknight1/retail-isa-investment/investment-service/handler"
	"github.com/oliknight1/retail-isa-investment/investment-service/model"
)

type mockService struct {
	createInvestment           func(customerId string, fundId string, amount float64) (*model.Investment, error)
	getInvestmentById          func(string) (*model.Investment, error)
	getInvestmentsByCustomerId func(string) (*[]model.Investment, error)
}

func (m *mockService) CreateInvestment(customerId string, fundId string, amount float64) (*model.Investment, error) {
	return m.createInvestment(customerId, fundId, amount)
}
func (m *mockService) GetInvestmentById(id string) (*model.Investment, error) {
	return m.getInvestmentById(id)
}
func (m *mockService) GetInvestmentsByCustomerId(id string) (*[]model.Investment, error) {
	return m.getInvestmentsByCustomerId(id)
}
func TestCreateInvestment(t *testing.T) {
	customerId := "cust-123"
	fundId := "fund-456"
	amount := 100.0

	mockService := &mockService{
		createInvestment: func(customerId string, fundId string, amount float64) (*model.Investment, error) {
			return &model.Investment{
				Id:         "test-id",
				CustomerId: customerId,
				FundId:     fundId,
				Amount:     amount,
				Status:     "pending",
				CreatedAt:  time.Now(),
			}, nil
		},
	}

	handler := handler.New(mockService)

	reqBody := fmt.Sprintf(`{"customerId":"%s","fundId":"%s","amount":%f}`, customerId, fundId, amount)
	req := httptest.NewRequest(http.MethodPost, "/investments", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateInvestment(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	var inv model.Investment
	if err := json.NewDecoder(res.Body).Decode(&inv); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	if inv.CustomerId != customerId || inv.FundId != fundId || inv.Amount != amount {
		t.Fatalf("unexpected investment returned: %+v", inv)
	}
}

func TestCreateInvestmentFailures(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "missing customerId",
			body:         `{"fundId":"fund-456","amount":100}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing fundId",
			body:         `{"customerId":"cust-123","amount":100}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "zero amount",
			body:         `{"customerId":"cust-123","fundId":"fund-456","amount":0}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid json",
			body:         `{"customerId":"abc",`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := handler.New(&mockService{})

			req := httptest.NewRequest(http.MethodPost, "/investments", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.CreateInvestment(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedCode {
				t.Fatalf("expected status %d, got %d", tc.expectedCode, res.StatusCode)
			}
		})
	}
}
