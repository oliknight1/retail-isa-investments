package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/oliknight1/retail-isa-investment/customer-service/internal"
	"github.com/oliknight1/retail-isa-investment/customer-service/service"
	"go.uber.org/zap"
)

type CustomerHandler struct {
	Service service.CustomerService
	Logger  *zap.Logger
}

func New(service service.CustomerService, logger *zap.Logger) *CustomerHandler {
	return &CustomerHandler{service, logger}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("handling CreateCustomer request",
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
	)
	internal.CustomerRequests.WithLabelValues("/customer", "POST").Inc()
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internal.CustomerCreationFailures.WithLabelValues("decode_error").Inc()
		h.Logger.Warn("Error decoding JSON",
			zap.Error(err),
		)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		internal.CustomerCreationFailures.WithLabelValues("missing_name").Inc()
		error := errors.New("name required")
		h.Logger.Warn("Error missing name",
			zap.Error(error),
		)
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}

	h.Logger.Info("calling RegisterCustomer",
		zap.String("name", req.Name),
	)
	customer, err := h.Service.RegisterCustomer(req.Name)
	if err != nil {
		h.Logger.Error("RegisterCustomer failed",
			zap.Error(err),
		)
		internal.CustomerCreationFailures.WithLabelValues("registration_failure").Inc()
		http.Error(w, "failed to register customer", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(customer); err != nil {
		h.Logger.Error("failed to encode JSON response",
			zap.Error(err),
			zap.String("customer_id", customer.Id),
		)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
	h.Logger.Info("customer registered successfully",
		zap.String("customer_id", customer.Id),
		zap.String("name", customer.Name),
	)
	internal.CustomerCreated.Inc()
}

func (h *CustomerHandler) GetCustomerById(w http.ResponseWriter, r *http.Request) {
	internal.CustomerRequests.WithLabelValues("/customer/{id}", "GET").Inc()
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	fmt.Println(parts)

	if len(parts) != 2 {
		internal.CustomerLookupFailures.WithLabelValues("invalid_url").Inc()
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	customerId := parts[1]
	if customerId == "" {
		internal.CustomerLookupFailures.WithLabelValues("missing_customer_id").Inc()
		h.Logger.Error("missing_customer_id")
		http.Error(w, "missing customer_id", http.StatusBadRequest)
		return
	}

	customer, err := h.Service.GetCustomerById(customerId)

	if err != nil {
		internal.CustomerLookupFailures.WithLabelValues("internal_server_error").Inc()
		h.Logger.Error("internal server error")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(customer); err != nil {
		h.Logger.Error("error encoding response", zap.Error(err))
		internal.CustomerLookupFailures.WithLabelValues("encoding_error").Inc()
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
