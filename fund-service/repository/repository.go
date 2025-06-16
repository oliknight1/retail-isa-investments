package repository

import (
	"encoding/json"
	"log"
	"os"

	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/model"
)

type Repository interface {
	GetFundById(id string) (*model.Fund, error)
	GetFundList() (*[]model.Fund, error)
}

type FundClient struct {
	Funds map[string]model.Fund
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
	fundMap := make(map[string]model.Fund)
	for _, fund := range fundList {
		fundMap[fund.Id] = fund
	}

	return &FundClient{Funds: fundMap}, nil
}

func (c *FundClient) GetFundById(id string) (*model.Fund, error) {
	fund, ok := c.Funds[id]
	if !ok {
		return nil, internal.FundNotFoundError(id)
	}
	return &fund, nil
}
