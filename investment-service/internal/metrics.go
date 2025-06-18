package internal

import "github.com/prometheus/client_golang/prometheus"

var (
	InvestmentRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "investment_requests_total",
			Help: "Total number of HTTP requests to investment service",
		},
		[]string{"endpoint", "method"},
	)

	InvestmentCreationFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "investment_creation_failures_total",
			Help: "Total number of failed investment creations",
		},
		[]string{"reason"},
	)
	InvestmentCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "investment_created_total",
			Help: "Total number of investments created successfully",
		},
	)
	InvestmentValidationEvents = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "investment_validation_events_total",
			Help: "Total number of investment validation events emitted",
		},
	)
)
