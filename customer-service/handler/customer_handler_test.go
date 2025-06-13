package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oliknight1/retail-isa-investment/customer-service/handler"
	"github.com/oliknight1/retail-isa-investment/customer-service/model"
)

type mockService struct {
	registerFn func(name string) (model.Customer, error)
}

func (s *mockService) RegisterCustomer(name string) (model.Customer, error) {
	return s.registerFn(name)
}

func TestCreateCustomerSuccess(t *testing.T) {
	expectedName := "Oli"
	mockService := &mockService{
		registerFn: func(name string) (model.Customer, error) {
			if name != expectedName {
				t.Fatalf("expected name %s, recieved %s", expectedName, name)
			}
			return model.Customer{Id: "1234", Name: name}, nil
		},
	}
	handler := &handler.CustomerHandler{Service: mockService}

	reqBody := fmt.Sprintf(`{"name":"%s"}`, expectedName)
	req := httptest.NewRequest(http.MethodPost, "/customers", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.CreateCustomer(recorder, req)

	expectedStatus := http.StatusCreated

	if status := recorder.Code; status != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, status)
	}

	var resp model.Customer
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Name != expectedName {
		t.Errorf("expected name %s, got %s", expectedName, resp.Name)
	}
}

func TestInvalidJSON(t *testing.T) {
	mockService := &mockService{
		registerFn: func(name string) (model.Customer, error) {
			return model.Customer{Id: "1234", Name: name}, nil
		},
	}
	handler := &handler.CustomerHandler{Service: mockService}

	reqBody := `{"name":"Oli"`
	req := httptest.NewRequest(http.MethodPost, "/customers", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.CreateCustomer(recorder, req)

	expectedStatus := http.StatusBadRequest

	if status := recorder.Code; status != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, status)
	}
}

func TestInavlidPayload(t *testing.T) {
	expectedName := "Oli"
	mockService := &mockService{
		registerFn: func(name string) (model.Customer, error) {
			if name != expectedName {
				t.Fatalf("expected name %s, recieved %s", expectedName, name)
			}
			return model.Customer{Id: "1234", Name: name}, nil
		},
	}
	handler := &handler.CustomerHandler{Service: mockService}

	reqBody := `{"name":""}`
	req := httptest.NewRequest(http.MethodPost, "/customers", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.CreateCustomer(recorder, req)

	expectedStatus := http.StatusBadRequest

	if status := recorder.Code; status != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, status)
	}
}

func TestInavlidService(t *testing.T) {
	mockService := &mockService{
		registerFn: func(name string) (model.Customer, error) {
			return model.Customer{}, errors.New("service error")
		},
	}
	handler := &handler.CustomerHandler{Service: mockService}

	reqBody := `{"name":"Oli"}`
	req := httptest.NewRequest(http.MethodPost, "/customers", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.CreateCustomer(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
}

// TODO: Test if content type if we always expect JSON
func TestInvalidContentType(t *testing.T) {}
