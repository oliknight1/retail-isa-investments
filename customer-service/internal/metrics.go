package internal

import "github.com/prometheus/client_golang/prometheus"

var (
	CustomerRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "customer_requests_total",
			Help: "Total number of HTTP requests received by customer service",
		},
		[]string{"endpoint", "method"},
	)
	CustomerCreationFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "customer_creation_failures_total",
			Help: "Total number of failed customer creation attempts",
		},
		[]string{"reason"},
	)
	CustomerCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "customer_created_total",
			Help: "Total number of successfully created customers",
		},
	)
	CustomerLookupFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "customer_lookup_failures_total",
			Help: "Total number of failed customer lookups",
		},
		[]string{"reason"},
	)
)
