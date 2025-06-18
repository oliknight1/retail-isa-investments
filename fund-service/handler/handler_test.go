package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/fund-service/handler"
	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/logger"
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
	logger := logger.NewMockLogger()
	handler := handler.FundHandler{Service: mockService, Logger: logger}

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
	logger := logger.NewMockLogger()
	handler := handler.FundHandler{Service: mockService, Logger: logger}

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

func TestGetFundByIdINotFound(t *testing.T) {
	mockService := &mockService{
		getFundById: func(id string) (*model.Fund, error) {
			return &model.Fund{}, internal.FundNotFoundError(id)
		},
	}

	logger := logger.NewMockLogger()
	handler := handler.FundHandler{Service: mockService, Logger: logger}

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

func TestGetFundListHandler(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		mockReturn     *[]model.Fund
		mockError      error
		expectedStatus int
		expectedBody   string // optional basic body check
	}{
		{
			name:           "returns funds successfully with valid riskLevel",
			queryParam:     "?riskLevel=Medium",
			mockReturn:     &[]model.Fund{{Id: "1", Name: "Fund 1", RiskLevel: "Medium"}},
			expectedStatus: http.StatusOK,
			expectedBody:   `"riskLevel":"Medium"`,
		},
		{
			name:           "returns all funds when riskLevel is not provided",
			queryParam:     "",
			mockReturn:     &[]model.Fund{{Id: "1", Name: "Fund 1", RiskLevel: "Low"}},
			expectedStatus: http.StatusOK,
			expectedBody:   `"riskLevel":"Low"`,
		},
		{
			name:           "returns 400 on invalid riskLevel",
			queryParam:     "?riskLevel=invalid",
			mockError:      internal.ErrInvalidRisklevel,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   internal.ErrInvalidRisklevel.Error(),
		},
		{
			name:           "returns 500 on unexpected error",
			queryParam:     "?riskLevel=Medium",
			mockError:      errors.New("something broke"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockService{
				getFundList: func(riskLevel *string) (*[]model.Fund, error) {
					return tt.mockReturn, tt.mockError
				},
			}

			logger := logger.NewMockLogger()
			handler := handler.FundHandler{Service: mockService, Logger: logger}
			req := httptest.NewRequest(http.MethodGet, "/funds"+tt.queryParam, nil)
			recorder := httptest.NewRecorder()

			handler.GetFundList(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			body := recorder.Body.String()
			if tt.expectedBody != "" && !strings.Contains(body, tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}
