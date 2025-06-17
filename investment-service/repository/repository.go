package repository

import (
	"fmt"
	"sync"

	"github.com/oliknight1/retail-isa-investment/investment-service/model"
)

type Repository interface {
	CreateInvestment(investment *model.Investment) error
	GetInvestmentById(id string) (*model.Investment, error)
	GetInvestmentsByCustomerId(id string) (*[]model.Investment, error)
}

type InvestmentClient struct {
	Investments map[string]model.Investment
	mu          sync.Mutex
}

func NewInvestmentClient() *InvestmentClient {
	return &InvestmentClient{
		Investments: make(map[string]model.Investment),
	}
}

func (c *InvestmentClient) CreateInvestment(investment model.Investment) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Investments[investment.Id] = investment
	return nil
}
func (c *InvestmentClient) GetInvestmentById(id string) (*model.Investment, error) {
	investment, ok := c.Investments[id]
	if !ok {
		return nil, fmt.Errorf("investment with id %s not found", id)
	}
	return &investment, nil

}
func (c *InvestmentClient) GetInvestmentsByCustomerId(id string) (*[]model.Investment, error) {
	var foundInvestments []model.Investment

	for _, investment := range c.Investments {
		if investment.CustomerId == id {
			foundInvestments = append(foundInvestments, investment)
		}
	}

	return &foundInvestments, nil
}
