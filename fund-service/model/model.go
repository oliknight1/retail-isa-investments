package model

type Fund struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FundAccount struct {
	CustomerID     string
	Balance        int64
	ReservedAmount int64
}
