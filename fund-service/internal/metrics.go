package internal

import "github.com/prometheus/client_golang/prometheus"

var (
	FundRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fund_requests_total",
			Help: "Total number of requests to fund endpoints",
		},
		[]string{"endpoint", "method"},
	)
	FundLookupFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fund_lookup_failures_total",
			Help: "Total number of failed fund lookups",
		},
		[]string{"reason"},
	)

	FundLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "fund_response_duration_seconds",
			Help:    "Response time for fund endpoints",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)
