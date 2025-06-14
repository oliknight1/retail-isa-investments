package service

import (
	"errors"

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

func (s *FundServiceImpl) GetFundList() (*[]model.Fund, error) {
	return s.repo.GetFundList()
}
