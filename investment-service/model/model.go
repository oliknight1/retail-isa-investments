package model

import "time"

type Investment struct {
	Id         string
	CustomerId string
	FundId     string
	Amount     float64
	// "pending", "completed", "failed"
	Status        string
	CreatedAt     time.Time
	CompletedAt   *time.Time
	FailureReason *string
}
