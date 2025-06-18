package handler_test

import (
	"encoding/json"
	"errors"
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

	handler := handler.New(mockService, nil)

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
			h := handler.New(&mockService{}, nil)

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

func TestGetInvestmentByIdSuccess(t *testing.T) {
	mockSvc := &mockService{
		getInvestmentById: func(id string) (*model.Investment, error) {
			return &model.Investment{
				Id:         id,
				CustomerId: "cust-123",
				FundId:     "fund-456",
				Amount:     100.0,
				Status:     "pending",
				CreatedAt:  time.Now(),
			}, nil
		},
	}
	handler := handler.New(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/investments/inv-123", nil)
	w := httptest.NewRecorder()

	handler.GetInvestmentById(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.StatusCode)
	}
	var inv model.Investment
	if err := json.NewDecoder(res.Body).Decode(&inv); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if inv.Id != "inv-123" {
		t.Errorf("expected id inv-123, got %s", inv.Id)
	}
}

func TestGetInvestmentByIdMissingId(t *testing.T) {
	mockSvc := &mockService{}
	handler := handler.New(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/investments/", nil)
	w := httptest.NewRecorder()

	handler.GetInvestmentById(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", res.StatusCode)
	}
}

func TestGetInvestmentByIdServiceError(t *testing.T) {
	mockSvc := &mockService{
		getInvestmentById: func(id string) (*model.Investment, error) {
			return nil, errors.New("db failure")
		},
	}
	handler := handler.New(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/investments/inv-123", nil)
	w := httptest.NewRecorder()

	handler.GetInvestmentById(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", res.StatusCode)
	}
}
func TestGetInvestmentsByCustomerIdSuccess(t *testing.T) {
	mockSvc := &mockService{
		getInvestmentsByCustomerId: func(id string) (*[]model.Investment, error) {
			return &[]model.Investment{
				{
					Id:         "inv-1",
					CustomerId: id,
					FundId:     "fund-1",
					Amount:     100.0,
					Status:     "pending",
					CreatedAt:  time.Now(),
				},
			}, nil
		},
	}
	handler := handler.New(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/customers/cust-123/investments", nil)
	w := httptest.NewRecorder()

	handler.GetInvestmentsByCustomerId(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.StatusCode)
	}
	var invs []model.Investment
	if err := json.NewDecoder(res.Body).Decode(&invs); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(invs) != 1 || invs[0].CustomerId != "cust-123" {
		t.Errorf("unexpected result: %+v", invs)
	}
}

func TestGetInvestmentsByCustomerIdMissingId(t *testing.T) {
	mockSvc := &mockService{}
	handler := handler.New(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/customers//investments", nil)
	w := httptest.NewRecorder()

	handler.GetInvestmentsByCustomerId(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", res.StatusCode)
	}
}

func TestGetInvestmentsByCustomerIdServiceError(t *testing.T) {
	mockSvc := &mockService{
		getInvestmentsByCustomerId: func(id string) (*[]model.Investment, error) {
			return nil, errors.New("some db error")
		},
	}
	handler := handler.New(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/customers/cust-123/investments", nil)
	w := httptest.NewRecorder()

	handler.GetInvestmentsByCustomerId(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", res.StatusCode)
	}
}
