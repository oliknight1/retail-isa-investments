package service

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/investment-service/event"
	"github.com/oliknight1/retail-isa-investment/investment-service/internal"
	"github.com/oliknight1/retail-isa-investment/investment-service/model"
	"github.com/oliknight1/retail-isa-investment/investment-service/repository"
)

type InvestmentService interface {
	CreateInvestment(customerId string, fundId string, amount float64) (*model.Investment, error)
	GetInvestmentById(string) (*model.Investment, error)
	GetInvestmentsByCustomerId(string) (*[]model.Investment, error)
}

type InvestmentServiceImpl struct {
	repo      repository.Repository
	publisher event.EventHandler
}

func New(repo repository.Repository, publisher event.EventHandler) *InvestmentServiceImpl {
	return &InvestmentServiceImpl{
		repo,
		publisher,
	}
}

func (s *InvestmentServiceImpl) CreateInvestment(customerId string, fundId string, amount float64) (*model.Investment, error) {
	if customerId == "" {
		return nil, internal.ErrMissingCustomerId
	}
	if fundId == "" {
		return nil, internal.ErrMissingFundId
	}
	if amount <= 0 {
		return nil, internal.ErrZeroTransactionAmount
	}
	investment := model.Investment{
		Id:         uuid.New().String(),
		CustomerId: customerId,
		FundId:     fundId,
		Amount:     amount,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}
	if err := s.repo.CreateInvestment(&investment); err != nil {
		return nil, err
	}

	s.publisher.Publish("investment.created", investment)

	if err := s.publisher.Publish("investment.processed", investment); err != nil {
		log.Println("error publishing investment.processed event: %v", err)
	}

	if err := s.publisher.Publish("investment.validation.pending", investment); err != nil {
		log.Println("error publishing investment.validation.pending event: %v", err)
	}

	return &investment, nil
}

func (s *InvestmentServiceImpl) GetInvestmentById(id string) (*model.Investment, error) {
	if id == "" {
		return nil, internal.ErrMissingFundId
	}
	return s.repo.GetInvestmentById(id)
}
func (s *InvestmentServiceImpl) GetInvestmentsByCustomerId(id string) (*[]model.Investment, error) {
	if id == "" {
		return nil, internal.ErrMissingCustomerId
	}
	return s.repo.GetInvestmentsByCustomerId(id)
}
