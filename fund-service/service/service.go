package service

import (
	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/logger"
	"github.com/oliknight1/retail-isa-investment/fund-service/model"
	"github.com/oliknight1/retail-isa-investment/fund-service/repository"
	"go.uber.org/zap"
)

type FundService interface {
	GetFundById(string) (*model.Fund, error)
	GetFundList(*string) (*[]model.Fund, error)
}

type FundServiceImpl struct {
	repo   repository.Repository
	Logger logger.Logger
}

func New(repo repository.Repository, logger logger.Logger) *FundServiceImpl {
	return &FundServiceImpl{
		repo,
		logger,
	}
}

func (s *FundServiceImpl) GetFundById(id string) (*model.Fund, error) {
	if id == "" {
		s.Logger.Error("missing fund_id when fetching fund", zap.Error(internal.ErrMissingId))
		return nil, internal.ErrMissingId
	}
	return s.repo.GetFundById(id)
}

func (s *FundServiceImpl) GetFundList(riskLevel *string) (*[]model.Fund, error) {
	allFunds, err := s.repo.GetFundList()
	if err != nil {
		s.Logger.Error("error fetching fund list", zap.Error(err))
		return nil, err
	}
	if riskLevel == nil {
		return allFunds, nil
	}

	//NOTE: This should be fetched from another service in real-app
	var riskOrder = map[string]int{
		"Low":    1,
		"Medium": 2,
		"High":   3,
	}

	allowedRisk, ok := riskOrder[*riskLevel]

	if !ok {
		s.Logger.Error("invalid riskLevel provided", riskLevel)
		return nil, internal.ErrInvalidRisklevel
	}
	appropiateFunds := []model.Fund{}

	//NOTE: allow for all risk levels at or below user risk level
	for _, fund := range *allFunds {
		if riskOrder[fund.RiskLevel] <= allowedRisk {
			appropiateFunds = append(appropiateFunds, fund)
		}
	}

	return &appropiateFunds, nil
}
