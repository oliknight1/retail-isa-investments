package service

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/oliknight1/retail-isa-investment/investment-service/event"
	"github.com/oliknight1/retail-isa-investment/investment-service/internal"
	"github.com/oliknight1/retail-isa-investment/investment-service/logger"
	"github.com/oliknight1/retail-isa-investment/investment-service/model"
	"github.com/oliknight1/retail-isa-investment/investment-service/repository"
	"go.uber.org/zap"
)

type InvestmentService interface {
	CreateInvestment(customerId string, fundId string, amount float64) (*model.Investment, error)
	GetInvestmentById(string) (*model.Investment, error)
	GetInvestmentsByCustomerId(string) (*[]model.Investment, error)
}

type InvestmentServiceImpl struct {
	repo      repository.Repository
	publisher event.EventHandler
	Logger    logger.Logger
}

func New(repo repository.Repository, publisher event.EventHandler, logger logger.Logger) *InvestmentServiceImpl {
	return &InvestmentServiceImpl{
		repo,
		publisher,
		logger,
	}
}

func (s *InvestmentServiceImpl) CreateInvestment(customerId string, fundId string, amount float64) (*model.Investment, error) {
	if customerId == "" {
		s.Logger.Error("missing customer_id in creation request", internal.ErrMissingCustomerId)
		return nil, internal.ErrMissingCustomerId
	}
	if fundId == "" {
		s.Logger.Error("missing fund_id in creation request", internal.ErrMissingFundId)
		return nil, internal.ErrMissingFundId
	}
	if amount <= 0 {
		s.Logger.Error("invalid transaction amount in creation request", internal.ErrZeroTransactionAmount)
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
	if err := s.repo.CreateInvestment(investment); err != nil {
		return nil, err
	}

	s.publisher.Publish("investment.created", investment)

	if err := s.publisher.Publish("investment.processed", investment); err != nil {
		s.Logger.Error("error publishing investment.processed event", zap.Error(err))
	}

	if err := s.publisher.Publish("investment.validation.pending", investment); err != nil {
		log.Println("error publishing investment.validation.pending event: %v", err)
		s.Logger.Error("error publishing investment.pending event", zap.Error(err))
	}
	internal.InvestmentValidationEvents.Inc()

	return &investment, nil
}

func (s *InvestmentServiceImpl) GetInvestmentById(id string) (*model.Investment, error) {
	if id == "" {
		s.Logger.Error("missing fund_id when requesting investment", zap.Error(internal.ErrMissingFundId))
		return nil, internal.ErrMissingFundId
	}
	return s.repo.GetInvestmentById(id)
}
func (s *InvestmentServiceImpl) GetInvestmentsByCustomerId(id string) (*[]model.Investment, error) {
	if id == "" {
		s.Logger.Error("missing customer_id when requesting investment", zap.Error(internal.ErrMissingCustomerId))
		return nil, internal.ErrMissingCustomerId
	}
	return s.repo.GetInvestmentsByCustomerId(id)
}
