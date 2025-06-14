package service_test

import (
	"testing"

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

func TestGetById(t *testing.T) {
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
	svc := service.New(mockRepo)

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

func TestGetList(t *testing.T) {
	expectedFundList := []model.Fund{
		{
			Id:          "fund-id",
			Name:        "fund-name",
			Description: "fund-description",
			RiskLevel:   "risk-level",
		},
		{
			Id:          "fund-id-2",
			Name:        "fund-name-2",
			Description: "fund-description-2",
			RiskLevel:   "risk-level",
		},
	}
	mockRepo := &mockRepo{
		getFundList: func() (*[]model.Fund, error) {
			return &expectedFundList, nil
		},
	}

	svc := service.New(mockRepo)

	funds, err := svc.GetFundList("risk-level")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(funds) == 0 {
		t.Errorf("expected funds length of : %d, got %d", len(expectedFundList), len(funds))
	}
	if len(funds) != len(expectedFundList) {
		t.Errorf("got %+v, want %+v", funds, expectedFundList)
	}
}
