package handler

import (
	"bytes"
	"encoding/json"
	"log"
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

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(customer); err != nil {
		log.Printf("failed to encode JSON: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}
