package handler

import (
	"bytes"
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

func writeJson(w http.ResponseWriter, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		log.Printf("failed to encode JSON: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
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
		if errors.Is(err, internal.ErrFundNotFound) {
			http.Error(w, internal.FundNotFoundError(fundId).Error(), http.StatusNotFound)
			return
		} else {
			fmt.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, fund)
}

func (h *FundHandler) GetFundList(w http.ResponseWriter, r *http.Request) {
	var riskLevel *string
	risk := r.URL.Query().Get("riskLevel")
	if risk != "" {
		riskLevel = &risk
	}

	funds, err := h.Service.GetFundList(riskLevel)
	if err != nil {
		if errors.Is(err, internal.ErrInvalidRisklevel) {
			http.Error(w, internal.ErrInvalidRisklevel.Error(), http.StatusBadRequest)
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	writeJson(w, funds)

}
