package repository_test

import (
	"testing"

	"github.com/oliknight1/retail-isa-investment/fund-service/model"
	"github.com/oliknight1/retail-isa-investment/fund-service/repository"
)

func TestGetById(t *testing.T) {
	fundId := "fund-ftse-100"
	expectedFund := model.Fund{
		Id:          fundId,
		Name:        "",
		Description: "",
		RiskLevel:   "",
	}

	db := &repository.InMemDb{
		Store: map[string]model.Fund{
			fundId: expectedFund,
		},
	}

	fund, err := db.GetFundById(fundId)

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
