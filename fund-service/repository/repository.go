package repository

import (
	"encoding/json"
	"log"
	"os"

	"github.com/oliknight1/retail-isa-investment/fund-service/model"
)

type Repository interface {
	GetFundById(id string) (*model.Fund, error)
	GetFundList() (*[]model.Fund, error)
}

type FundClient struct {
	Funds []model.Fund
}

// NOTE: replace with a client that fetches / listens to external service providing the fund list
// Currently fetches from local file
func NewFundClient(path string) (*FundClient, error) {
	file, err := os.Open(path)

	if err != nil {
		log.Printf("error reading from file: %v", err)
		return nil, err
	}
	defer file.Close()

	var fundList []model.Fund
	if err := json.NewDecoder(file).Decode(&fundList); err != nil {
		log.Printf("error deconding fund data: %v", err)
		return nil, err
	}

	return &FundClient{Funds: fundList}, nil
}

func (c *FundClient) GetFundById(id string) (*model.Fund, error) {
	var foundFund *model.Fund
	for _, fund := range c.Funds {
		if fund.Id == id {
			foundFund = &fund
			break
		}
	}
	return foundFund, nil
}
func (c *FundClient) GetFundList(riskLevel *string) (*[]model.Fund, error) {
	return &c.Funds, nil
}
