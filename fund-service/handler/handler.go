package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/service"
)

type FundHandler struct {
	Service service.FundService
}

func New(service service.FundService) *FundHandler {
	return &FundHandler{service}
}

func (h *FundHandler) GetFundById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 2 {
		http.Error(w, internal.ErrInvalidUrl.Error(), http.StatusBadRequest)
		return
	}

	fundId := parts[1]
	if fundId == "" {
		http.Error(w, internal.ErrMissingId.Error(), http.StatusBadRequest)
		return
	}

	fund, err := h.Service.GetFundById(fundId)

	if err != nil {
		switch {
		case errors.Is(err, internal.ErrInvalidId):
			http.Error(w, internal.ErrInvalidId.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, internal.ErrFundNotFound):
			http.Error(w, internal.FundNotFoundError(fundId).Error(), http.StatusNotFound)
			return
		default:
			fmt.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(fund); err != nil {
		log.Printf("failed to encode JSON: %v", err)
	}
}
