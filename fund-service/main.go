package main

import (
	"log"
	"net/http"

	"github.com/oliknight1/retail-isa-investment/fund-service/handler"
	"github.com/oliknight1/retail-isa-investment/fund-service/repository"
	"github.com/oliknight1/retail-isa-investment/fund-service/service"
)

func main() {
	repo, err := repository.NewFundClient("./repository/funds.json")

	if err != nil {

		log.Fatalf("Error reading funds.json: \n%v\n", err)
	}
	svc := service.New(repo)
	fh := handler.New(svc)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	http.HandleFunc("/funds", fh.GetFundList)

	http.HandleFunc("/funds/", fh.GetFundById)

	log.Println("fund-service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
