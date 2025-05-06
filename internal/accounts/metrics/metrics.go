package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "accounts_requests_total",
			Help: "Total HTTP requests for Accounts service.",
		},
		[]string{"method", "route", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "accounts_request_duration_seconds",
			Help:    "HTTP request latency distributions.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestsTotal, RequestDuration)
}
