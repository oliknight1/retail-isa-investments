package service

import (
	"log"

	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/model"
	"github.com/oliknight1/retail-isa-investment/fund-service/repository"
)

type FundService interface {
	GetFundById(string) (*model.Fund, error)
	GetFundList(*string) (*[]model.Fund, error)
}

type FundServiceImpl struct {
	repo repository.Repository
}

func New(repo repository.Repository) *FundServiceImpl {
	return &FundServiceImpl{
		repo,
	}
}

func (s *FundServiceImpl) GetFundById(id string) (*model.Fund, error) {
	if id == "" {
		return nil, internal.ErrMissingId
	}
	if err := uuid.Validate(id); err != nil {
		log.Printf("invalid UUID provided: %s, error: %v", id, err)
		return nil, internal.ErrInvalidId
	}
	return s.repo.GetFundById(id)
}

func (s *FundServiceImpl) GetFundList(riskLevel *string) ([]model.Fund, error) {
	allFunds, err := s.repo.GetFundList()
	if err != nil {
		return nil, err
	}
	if riskLevel == nil {
		return *allFunds, nil
	}

	//NOTE: This should be fetched from another service in real-app
	var riskOrder = map[string]int{
		"Low":    1,
		"Medium": 2,
		"High":   3,
	}

	allowedRisk, ok := riskOrder[*riskLevel]

	if !ok {
		return nil, internal.ErrInvalidRisklevel
	}
	appropiateFunds := []model.Fund{}

	//NOTE: allow for all risk levels at or below user risk level
	for _, fund := range *allFunds {
		if riskOrder[fund.RiskLevel] <= allowedRisk {
			appropiateFunds = append(appropiateFunds, fund)
		}
	}

	return appropiateFunds, nil
}
