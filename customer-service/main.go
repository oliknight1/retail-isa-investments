package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oliknight1/retail-isa-investment/customer-service/event"
	"github.com/oliknight1/retail-isa-investment/customer-service/handler"
	"github.com/oliknight1/retail-isa-investment/customer-service/repository"
	"github.com/oliknight1/retail-isa-investment/customer-service/service"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	repo := repository.New()
	pub := event.NewNatsPublisher(natsURL)
	svc := service.New(repo, pub)
	ch := handler.New(svc)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("POST /customers", ch.CreateCustomer)

	log.Println("Customer service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
