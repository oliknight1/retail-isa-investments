package internal

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidId        = errors.New("invalid UUID provided")
	ErrMissingId        = errors.New("id is required")
	ErrInvalidRisklevel = errors.New("invalid risk level")
	ErrInvalidUrl       = errors.New("invalid url")
	ErrFundNotFound     = errors.New("fund not found")
)

func FundNotFoundError(id string) error {
	return fmt.Errorf("%w: %s", ErrFundNotFound, id)
}
