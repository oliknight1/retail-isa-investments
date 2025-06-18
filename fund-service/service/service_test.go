package service_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/model"
	"github.com/oliknight1/retail-isa-investment/fund-service/service"
)

type mockRepo struct {
	getFundByIdFn func(id string) (*model.Fund, error)
	getFundList   func() (*[]model.Fund, error)
}

func (m *mockRepo) GetFundById(id string) (*model.Fund, error) {
	return m.getFundByIdFn(id)
}
func (m *mockRepo) GetFundList() (*[]model.Fund, error) {
	return m.getFundList()
}

func TestGetByIdSuccess(t *testing.T) {
	expectedFund := model.Fund{
		Id:          "fund-id",
		Name:        "fund-name",
		Description: "fund-description",
		RiskLevel:   "risk-level",
	}
	mockRepo := &mockRepo{
		getFundByIdFn: func(id string) (*model.Fund, error) {
			if id != expectedFund.Id {
				t.Fatalf("expected ID %s, got %s", expectedFund.Id, id)
			}
			return &expectedFund, nil
		},
	}
	svc := service.New(mockRepo, nil)

	fund, err := svc.GetFundById(expectedFund.Id)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if fund == nil {
		t.Errorf("expected fund: %s, got nil", expectedFund.Name)
	}
	if *fund != expectedFund {
		t.Errorf("got %+v, want %+v", *fund, expectedFund)
	}

}

func TestGetFundByIdEmptyID(t *testing.T) {
	mockRepo := &mockRepo{
		getFundByIdFn: func(id string) (*model.Fund, error) {
			t.Fatal("repo.GetFundById should NOT be called for empty ID")
			return nil, nil
		},
	}
	svc := service.New(mockRepo, nil)

	_, err := svc.GetFundById("")

	if err == nil {
		t.Error("expected error for empty fund id")
		return
	}

	if err.Error() != internal.ErrMissingId.Error() {
		t.Errorf("expected '%s' error, got: %v", internal.ErrMissingId, err)
	}
}
func TestGetFundByIdNotFound(t *testing.T) {
	mockRepo := &mockRepo{
		getFundByIdFn: func(id string) (*model.Fund, error) {
			return nil, fmt.Errorf("fund with id: %s not found", id)
		},
	}
	svc := service.New(mockRepo, nil)

	_, err := svc.GetFundById("fund-doesn't-exist")

	if err == nil {
		t.Error("expected error for fund not found")
		return
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestGetListSuccess(t *testing.T) {
	expectedFundList := []model.Fund{
		{
			Id:          "fund-id",
			Name:        "fund-name",
			Description: "fund-description",
			RiskLevel:   "Low",
		},
		{
			Id:          "fund-id-2",
			Name:        "fund-name-2",
			Description: "fund-description-2",
			RiskLevel:   "Low",
		},
	}
	mockRepo := &mockRepo{
		getFundList: func() (*[]model.Fund, error) {
			return &expectedFundList, nil
		},
	}

	svc := service.New(mockRepo, nil)

	riskLevel := "Low"
	funds, err := svc.GetFundList(&riskLevel)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(*funds) == 0 {
		t.Errorf("expected funds length of : %d, got %d", len(expectedFundList), len(*funds))
	}
	//TODO: proper equivilance test
	if len(*funds) != len(expectedFundList) {
		t.Errorf("got %+v, want %+v", funds, expectedFundList)
	}
}
func TestGetFundListFiltersRiskLevel(t *testing.T) {
	// risk levels
	var (
		low     = "Low"
		medium  = "Medium"
		high    = "High"
		invalid = "low"
	)
	tests := []struct {
		name             string
		riskLevel        *string
		fundList         []model.Fund
		expectedFundList []model.Fund
		expectErr        bool
	}{
		{
			name:      "filters out funds above Medium",
			riskLevel: &medium,
			fundList: []model.Fund{
				{Id: "1", Name: "Fund 1", RiskLevel: "Low"},
				{Id: "2", Name: "Fund 2", RiskLevel: "Low"},
				{Id: "3", Name: "Fund 3", RiskLevel: "Medium"},
				{Id: "4", Name: "Fund 4", RiskLevel: "High"},
			},
			expectedFundList: []model.Fund{
				{Id: "1", Name: "Fund 1", RiskLevel: "Low"},
				{Id: "2", Name: "Fund 2", RiskLevel: "Low"},
				{Id: "3", Name: "Fund 3", RiskLevel: "Medium"},
			},
			expectErr: false,
		},
		{
			name:      "includes all when riskLevel is High",
			riskLevel: &high,
			fundList: []model.Fund{
				{Id: "1", Name: "Fund 1", RiskLevel: "Low"},
				{Id: "2", Name: "Fund 2", RiskLevel: "Medium"},
				{Id: "3", Name: "Fund 3", RiskLevel: "High"},
			},
			expectedFundList: []model.Fund{
				{Id: "1", Name: "Fund 1", RiskLevel: "Low"},
				{Id: "2", Name: "Fund 2", RiskLevel: "Medium"},
				{Id: "3", Name: "Fund 3", RiskLevel: "High"},
			},
			expectErr: false,
		},
		{
			name:      "filters out all when riskLevel is Low",
			riskLevel: &low,
			fundList: []model.Fund{
				{Id: "1", Name: "Fund 1", RiskLevel: "Medium"},
				{Id: "2", Name: "Fund 2", RiskLevel: "High"},
			},
			expectedFundList: []model.Fund{},
			expectErr:        false,
		},
		{
			name:      "invalid risk level: lowercase input",
			riskLevel: &invalid,
			fundList: []model.Fund{
				{Id: "1", Name: "Fund 1", RiskLevel: "Low"},
				{Id: "2", Name: "Fund 2", RiskLevel: "Medium"},
			},
			expectedFundList: nil,
			expectErr:        true,
		},
		{
			name:      "empty risk level",
			riskLevel: nil,
			fundList: []model.Fund{
				{Id: "1", Name: "Fund 1", RiskLevel: "Low"},
				{Id: "2", Name: "Fund 2", RiskLevel: "Medium"},
			},
			expectedFundList: []model.Fund{
				{Id: "1", Name: "Fund 1", RiskLevel: "Low"},
				{Id: "2", Name: "Fund 2", RiskLevel: "Medium"},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRepo{
				getFundList: func() (*[]model.Fund, error) {
					return &tt.fundList, nil
				},
			}

			svc := service.New(mock, nil)

			actual, err := svc.GetFundList(tt.riskLevel)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if diff := cmp.Diff(&tt.expectedFundList, actual); diff != "" {
				t.Errorf("unexpected fund list (-want +got):\n%s", diff)
			}
		})
	}
}
