package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/oliknight1/retail-isa-investment/investment-service/event"
	"github.com/oliknight1/retail-isa-investment/investment-service/handler"
	"github.com/oliknight1/retail-isa-investment/investment-service/internal"
	"github.com/oliknight1/retail-isa-investment/investment-service/logger"
	"github.com/oliknight1/retail-isa-investment/investment-service/repository"
	"github.com/oliknight1/retail-isa-investment/investment-service/service"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	prometheus.MustRegister(
		internal.InvestmentRequests,
		internal.InvestmentCreated,
		internal.InvestmentCreationFailures,
		internal.InvestmentValidationEvents,
	)

	repo := repository.NewInvestmentClient()
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	publisher, err := event.NewNatsPublisher(natsURL)
	if err != nil {
		log.Printf("error connecting to publisher: %v", err)
	}
	svc := service.New(repo, publisher, logger)
	ih := handler.New(svc, logger)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	http.HandleFunc("POST /investments", ih.CreateInvestment)

	http.HandleFunc("GET /investments/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/investments/customer/") {
			fmt.Println("PREEFIX")
			ih.GetInvestmentsByCustomerId(w, r)
			return
		}
		ih.GetInvestmentById(w, r)
	})

	log.Println("Customer service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
