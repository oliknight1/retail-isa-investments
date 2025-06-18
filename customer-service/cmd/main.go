package main

import (
	"log"
	"net/http"
	"os"

	"github.com/oliknight1/retail-isa-investment/customer-service/event"
	"github.com/oliknight1/retail-isa-investment/customer-service/handler"
	"github.com/oliknight1/retail-isa-investment/customer-service/internal"
	"github.com/oliknight1/retail-isa-investment/customer-service/repository"
	"github.com/oliknight1/retail-isa-investment/customer-service/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	prometheus.MustRegister(
		internal.CustomerCreated,
		internal.CustomerCreationFailures,
		internal.CustomerRequests,
		internal.CustomerLookupFailures,
	)

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	repo := repository.New()
	pub := event.NewNatsPublisher(natsURL)
	svc := service.New(repo, pub)
	ch := handler.New(svc, logger)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("POST /customer", ch.CreateCustomer)

	http.HandleFunc("GET /customer/", ch.GetCustomerById)

	logger.Info("Customer service running on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}

}
