package repository_test

import (
	"testing"
	"time"

	"github.com/oliknight1/retail-isa-investment/investment-service/model"
	"github.com/oliknight1/retail-isa-investment/investment-service/repository"
)

func initDb(db *repository.InvestmentClient) {
	now := time.Now()
	later := now.Add(1 * time.Hour)
	failReason := "insufficient funds"

	investments := []model.Investment{
		{
			Id:          "inv-1",
			CustomerId:  "cust-1",
			FundId:      "fund-1",
			Amount:      1000.00,
			Status:      "completed",
			CreatedAt:   now.Add(-48 * time.Hour),
			CompletedAt: &now,
		},
		{
			Id:          "inv-2",
			CustomerId:  "cust-1",
			FundId:      "fund-2",
			Amount:      500.00,
			Status:      "pending",
			CreatedAt:   now.Add(-24 * time.Hour),
			CompletedAt: nil,
		},
		{
			Id:          "inv-3",
			CustomerId:  "cust-2",
			FundId:      "fund-1",
			Amount:      2000.00,
			Status:      "completed",
			CreatedAt:   now.Add(-72 * time.Hour),
			CompletedAt: &later,
		},
		{
			Id:            "inv-4",
			CustomerId:    "cust-3",
			FundId:        "fund-3",
			Amount:        750.00,
			Status:        "failed",
			CreatedAt:     now.Add(-12 * time.Hour),
			FailureReason: &failReason,
		},
		{
			Id:          "inv-5",
			CustomerId:  "cust-1",
			FundId:      "fund-3",
			Amount:      300.00,
			Status:      "completed",
			CreatedAt:   now.Add(-6 * time.Hour),
			CompletedAt: &now,
		},
	}

	for _, investment := range investments {
		db.CreateInvestment(investment)
	}

}

func TestGetByCustomerId(t *testing.T) {
	expectedId := "cust-1"
	db := repository.NewInvestmentClient()
	initDb(db)

	investments, err := db.GetInvestmentsByCustomerId(expectedId)

	if err != nil {
		t.Fatalf("unexpect error: %v", err)
	}
	if len(*investments) != 3 {
		t.Fatalf("expected 3 investments, recieved: %d", len(*investments))
	}

	for _, investment := range *investments {
		if investment.CustomerId != expectedId {
			t.Fatalf("expected id: %s, recieved: %s", expectedId, investment.CustomerId)
		}
	}

}
