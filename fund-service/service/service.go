package service

import (
	"errors"
	"fmt"

	"github.com/oliknight1/retail-isa-investment/fund-service/model"
	"github.com/oliknight1/retail-isa-investment/fund-service/repository"
)

type FundService interface {
	GetFundById(string) (*model.Fund, error)
	GetFundList() (*[]model.Fund, error)
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
		return nil, errors.New("fund id is required")
	}
	return s.repo.GetFundById(id)
}

func (s *FundServiceImpl) GetFundList(riskLevel string) ([]model.Fund, error) {
	//NOTE: This should be fetched from another service in real-app
	var riskOrder = map[string]int{
		"Low":    1,
		"Medium": 2,
		"High":   3,
	}
	allFunds, err := s.repo.GetFundList()
	if err != nil {
		return nil, err
	}

	allowedRisk, ok := riskOrder[riskLevel]

	if !ok {
		return nil, fmt.Errorf("invalid risk level: %s", riskLevel)
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
