package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/oliknight1/retail-isa-investment/investment-service/internal"
	"github.com/oliknight1/retail-isa-investment/investment-service/service"
)

type InvestmentHandler struct {
	Service service.InvestmentService
}

func New(service service.InvestmentService) *InvestmentHandler {
	return &InvestmentHandler{service}
}

func (h *InvestmentHandler) CreateInvestment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerId string  `json:"customerId"`
		FundId     string  `json:"fundId"`
		Amount     float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error reading JSON: \n%s", err), http.StatusBadRequest)
		return
	}

	if req.CustomerId == "" {
		http.Error(w, internal.ErrMissingCustomerId.Error(), http.StatusBadRequest)
		return
	}
	if req.FundId == "" {
		http.Error(w, internal.ErrMissingFundId.Error(), http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, internal.ErrZeroTransactionAmount.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.Service.CreateInvestment(req.CustomerId, req.FundId, req.Amount)

	if err != nil {
		http.Error(w, "failed to create investment", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(transaction); err != nil {
		log.Printf("failed to encode JSON: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func (h *InvestmentHandler) GetInvestmentById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, internal.ErrMissingFundId.Error(), http.StatusBadRequest)
		return
	}
	id := parts[2]

	investment, err := h.Service.GetInvestmentById(id)
	if err != nil {
		http.Error(w, "failed to get investment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(investment)
}

func (h *InvestmentHandler) GetInvestmentsByCustomerId(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, internal.ErrMissingCustomerId.Error(), http.StatusBadRequest)
		return
	}
	customerId := parts[3]

	investments, err := h.Service.GetInvestmentsByCustomerId(customerId)
	if err != nil {
		http.Error(w, "failed to get investments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(investments)
}
