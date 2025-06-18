package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/oliknight1/retail-isa-investment/investment-service/internal"
	"github.com/oliknight1/retail-isa-investment/investment-service/logger"
	"github.com/oliknight1/retail-isa-investment/investment-service/service"
	"go.uber.org/zap"
)

type InvestmentHandler struct {
	Service service.InvestmentService
	Logger  logger.Logger
}

func New(service service.InvestmentService, logger logger.Logger) *InvestmentHandler {
	return &InvestmentHandler{service, logger}
}

func (h *InvestmentHandler) CreateInvestment(w http.ResponseWriter, r *http.Request) {
	internal.InvestmentRequests.WithLabelValues("/investments", "POST").Inc()
	var req struct {
		CustomerId string  `json:"customerId"`
		FundId     string  `json:"fundId"`
		Amount     float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("failed to decode investment creation request", zap.Error(err))
		internal.InvestmentCreationFailures.WithLabelValues("decode_error").Inc()
		http.Error(w, fmt.Sprintf("error reading JSON: \n%s", err), http.StatusBadRequest)
		return
	}

	if req.CustomerId == "" {
		h.Logger.Error("missing customer_id in creation request", internal.ErrMissingCustomerId)
		internal.InvestmentCreationFailures.WithLabelValues("missing_customer_id").Inc()
		http.Error(w, internal.ErrMissingCustomerId.Error(), http.StatusBadRequest)
		return
	}
	if req.FundId == "" {
		h.Logger.Error("missing fund_id in creation request", internal.ErrMissingFundId)
		internal.InvestmentCreationFailures.WithLabelValues("mising_fund_id").Inc()
		http.Error(w, internal.ErrMissingFundId.Error(), http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		h.Logger.Error("invalid transaction amount in creation request", internal.ErrZeroTransactionAmount)
		internal.InvestmentCreationFailures.WithLabelValues("invalid_amount").Inc()
		http.Error(w, internal.ErrZeroTransactionAmount.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.Service.CreateInvestment(req.CustomerId, req.FundId, req.Amount)

	if err != nil {
		h.Logger.Error("internal service error", zap.Error(err))
		internal.InvestmentCreationFailures.WithLabelValues("service_error").Inc()
		http.Error(w, "failed to create investment", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(transaction); err != nil {
		h.Logger.Error("failed to write transaction creation to JSON", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	internal.InvestmentCreated.Inc()
	h.Logger.Info("transaction successfully created", transaction.Id)
	w.Write(buf.Bytes())
}

func (h *InvestmentHandler) GetInvestmentById(w http.ResponseWriter, r *http.Request) {
	internal.InvestmentRequests.WithLabelValues("/investments/{id}", "GET").Inc()
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		h.Logger.Error("missing fund_id when requesting investment", zap.Error(internal.ErrMissingFundId))
		http.Error(w, internal.ErrMissingFundId.Error(), http.StatusBadRequest)
		return
	}
	id := parts[2]

	investment, err := h.Service.GetInvestmentById(id)
	if err != nil {
		h.Logger.Error("failed to get investment", zap.Error(err))
		http.Error(w, "failed to get investment", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("investment found", investment.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(investment)
}

func (h *InvestmentHandler) GetInvestmentsByCustomerId(w http.ResponseWriter, r *http.Request) {
	internal.InvestmentRequests.WithLabelValues("/investments/customer/{id}", "GET").Inc()
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		h.Logger.Error("missing customer_id when requesting investment", zap.Error(internal.ErrMissingCustomerId))
		http.Error(w, internal.ErrMissingCustomerId.Error(), http.StatusBadRequest)
		return
	}
	customerId := parts[3]

	investments, err := h.Service.GetInvestmentsByCustomerId(customerId)
	if err != nil {
		h.Logger.Error("failed to get investment", zap.Error(err))
		http.Error(w, "failed to get investments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.Logger.Info("investment found", len(*investments))
	json.NewEncoder(w).Encode(investments)
}
