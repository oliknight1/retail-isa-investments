package internal

import "errors"

var (
	ErrMissingCustomerId     = errors.New("customer id is required")
	ErrMissingFundId         = errors.New("fund id is required")
	ErrZeroTransactionAmount = errors.New("transaction amount must be greater than 0")
)
