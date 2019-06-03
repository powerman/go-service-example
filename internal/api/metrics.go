package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metric struct { //nolint:gochecknoglobals
	reqInFlight prometheus.Gauge
	reqTotal    *prometheus.CounterVec
	reqDuration *prometheus.HistogramVec
}

const (
	resourceLabel = "resource"
	methodLabel   = "method"
	codeLabel     = "code"
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(namespace string) {
	const subsystem = "api"

	metric.reqInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_in_flight",
			Help:      "Amount of currently processing API requests.",
		},
	)
	metric.reqTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_total",
			Help:      "Amount of processed API requests.",
		},
		[]string{resourceLabel, methodLabel, codeLabel},
	)
	metric.reqDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_request_duration_seconds",
			Help:      "API request latency distributions.",
		},
		[]string{resourceLabel, methodLabel, codeLabel},
	)

	var (
		resources = [...]string{
			"/contacts",
		}
		methods = [...]string{
			"GET",
			"POST",
		}
		codes = [...]string{
			"200",
			"201",
			"400",
			"401",
			"403",
			"404",
			"405",
			"422",
			"500",
		}
	)
	for _, resource := range resources {
		for _, method := range methods {
			for _, code := range codes {
				l := prometheus.Labels{
					resourceLabel: resource,
					methodLabel:   method,
					codeLabel:     code,
				}
				metric.reqTotal.With(l)
				metric.reqDuration.With(l)
			}
		}
	}
}
