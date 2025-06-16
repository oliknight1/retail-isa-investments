package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/fund-service/handler"
	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/model"
)

type mockService struct {
	getFundById func(id string) (*model.Fund, error)
	getFundList func(riskLevel *string) (*[]model.Fund, error)
}

func (s *mockService) GetFundById(id string) (*model.Fund, error) {
	return s.getFundById(id)
}
func (s *mockService) GetFundList(riskLevel *string) (*[]model.Fund, error) {
	return s.getFundList(riskLevel)
}

func TestGetFundByIdSuccess(t *testing.T) {
	id := uuid.New().String()
	expectedFund := model.Fund{
		Id:          id,
		Name:        "fund",
		Description: "Fund Description",
		RiskLevel:   "Medium",
	}
	mockService := &mockService{
		getFundById: func(id string) (*model.Fund, error) {
			if id != expectedFund.Id {
				t.Fatalf("expected id %s, recieved %s", expectedFund.Id, id)
			}
			return &expectedFund, nil
		},
	}
	handler := handler.FundHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/funds/%s", id), nil)

	recorder := httptest.NewRecorder()

	handler.GetFundById(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}

	var result model.Fund
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Id != expectedFund.Id {
		t.Errorf("expected fund id %s, got %s", expectedFund.Id, result.Id)
	}

}

func TestGetFundByIdMissingId(t *testing.T) {
	mockService := &mockService{
		getFundById: func(id string) (*model.Fund, error) {
			return &model.Fund{}, internal.ErrMissingId
		},
	}
	handler := &handler.FundHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/funds/", nil)

	recorder := httptest.NewRecorder()

	handler.GetFundById(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, internal.ErrInvalidUrl.Error()) {
		t.Errorf("expected body to contain '%s', got '%s'", internal.ErrInvalidUrl, body)
	}

}
func TestGetFundByIdInvalidId(t *testing.T) {
	mockService := &mockService{
		getFundById: func(id string) (*model.Fund, error) {
			return &model.Fund{}, internal.ErrInvalidId
		},
	}
	handler := &handler.FundHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/funds/fund-123", nil)

	recorder := httptest.NewRecorder()

	handler.GetFundById(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", recorder.Code)
	}

	body := strings.TrimSpace(recorder.Body.String())
	if body != string(internal.ErrInvalidId.Error()) {
		t.Errorf("expected error message: %s got: %s", internal.ErrInvalidId, body)
	}

}
func TestGetFundByIdINotFound(t *testing.T) {
	mockService := &mockService{
		getFundById: func(id string) (*model.Fund, error) {
			return &model.Fund{}, internal.FundNotFoundError(id)
		},
	}

	handler := &handler.FundHandler{Service: mockService}

	id := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/funds/%s", id), nil)
	recorder := httptest.NewRecorder()

	handler.GetFundById(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", recorder.Code)
	}

	body := strings.TrimSpace(recorder.Body.String())
	if body != string(internal.FundNotFoundError(id).Error()) {
		t.Errorf("expected body to contain '%s', got '%s'", internal.FundNotFoundError(id), body)
	}
}
