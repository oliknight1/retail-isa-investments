package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/logger"
	"github.com/oliknight1/retail-isa-investment/fund-service/service"
	"go.uber.org/zap"
)

type FundHandler struct {
	Service service.FundService
	Logger  logger.Logger
}

func New(service service.FundService, logger logger.Logger) *FundHandler {
	return &FundHandler{service, logger}
}

func (h *FundHandler) writeJson(w http.ResponseWriter, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		log.Printf("failed to encode JSON: %v", err)
		h.Logger.Error("failed to endcode JSON", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (h *FundHandler) GetFundById(w http.ResponseWriter, r *http.Request) {
	internal.FundRequests.WithLabelValues("/funds/{id}", "GET").Inc()
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 2 {
		internal.FundLookupFailures.WithLabelValues("invalid_url")
		h.Logger.Error("invalid url", zap.Error(internal.ErrInvalidUrl))
		http.Error(w, internal.ErrInvalidUrl.Error(), http.StatusBadRequest)
		return
	}

	fundId := parts[1]
	if fundId == "" {
		internal.FundLookupFailures.WithLabelValues("missing_id")
		h.Logger.Error("missing fund_id", zap.Error(internal.ErrMissingId))
		http.Error(w, internal.ErrMissingId.Error(), http.StatusBadRequest)
		return
	}

	fund, err := h.Service.GetFundById(fundId)

	if err != nil {
		if errors.Is(err, internal.ErrFundNotFound) {
			internal.FundLookupFailures.WithLabelValues("not_found")
			h.Logger.Error("fund not found", zap.Error(internal.ErrFundNotFound))
			http.Error(w, internal.FundNotFoundError(fundId).Error(), http.StatusNotFound)
			return
		} else {
			internal.FundLookupFailures.WithLabelValues("internal_error")
			h.Logger.Error("internal server error")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	h.writeJson(w, fund)
	h.Logger.Info("successfully found fund", fund.Id)
}

func (h *FundHandler) GetFundList(w http.ResponseWriter, r *http.Request) {
	internal.FundRequests.WithLabelValues("/funds", r.Method).Inc()
	var riskLevel *string
	risk := r.URL.Query().Get("riskLevel")
	if risk != "" {
		riskLevel = &risk
	}

	funds, err := h.Service.GetFundList(riskLevel)
	if err != nil {
		if errors.Is(err, internal.ErrInvalidRisklevel) {
			h.Logger.Error("invalid riskLevel", zap.Error(internal.ErrInvalidRisklevel))
			http.Error(w, internal.ErrInvalidRisklevel.Error(), http.StatusBadRequest)
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		h.Logger.Error("internal server error")
	}
	h.writeJson(w, funds)

}
