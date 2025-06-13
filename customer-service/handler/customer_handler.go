package handler

import (
	"encoding/json"
	"net/http"

	"github.com/oliknight1/retail-isa-investment/customer-service/service"
)

type CustomerHandler struct {
	Service service.CustomerService
}

func New(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	customer, err := h.Service.RegisterCustomer(req.Name)
	if err != nil {
		http.Error(w, "failed to register customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}
