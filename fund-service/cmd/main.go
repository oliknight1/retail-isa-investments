package main

import (
	"log"
	"net/http"
	"os"

	"github.com/oliknight1/retail-isa-investment/fund-service/handler"
	"github.com/oliknight1/retail-isa-investment/fund-service/internal"
	"github.com/oliknight1/retail-isa-investment/fund-service/logger"
	"github.com/oliknight1/retail-isa-investment/fund-service/repository"
	"github.com/oliknight1/retail-isa-investment/fund-service/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()
	jsonFilePath := os.Getenv("FUNDS_JSON_PATH")
	if jsonFilePath == "" {
		jsonFilePath = "./repository/funds.json"
	}
	repo, err := repository.NewFundClient(jsonFilePath)
	prometheus.MustRegister(internal.FundLookupFailures, internal.FundRequests)

	if err != nil {
		logger.Error("Error reading funds.json", zap.Error(err))
	}
	svc := service.New(repo, logger)
	fh := handler.New(svc, logger)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/funds", fh.GetFundList)

	http.HandleFunc("/funds/", fh.GetFundById)

	logger.Info("fund-service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("server failed to start", zap.Error(err))
	}
}
