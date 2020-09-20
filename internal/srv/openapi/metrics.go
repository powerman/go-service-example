package openapi

import (
	"strconv"

	"github.com/go-openapi/loads"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-service-example/api/openapi/restapi"
)

// Metric contains general metrics for OpenAPI methods.
var metric struct { //nolint:gochecknoglobals // Metrics are global anyway.
	reqInFlight prometheus.Gauge
	reqTotal    *prometheus.CounterVec
	reqDuration *prometheus.HistogramVec
}

const (
	resourceLabel = "resource"
	methodLabel   = "method"
	codeLabel     = "code"
	failedLabel   = "failed"
)

//nolint:gochecknoglobals // By design.
var (
	// Initialized with codes returned by go-swagger and middlewares
	// after metrics middleware (accessLog).
	codeLabels = []int{400, 401, 403, 422}
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(reg *prometheus.Registry, namespace string) {
	const subsystem = "openapi"

	metric.reqInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_in_flight",
			Help:      "Amount of currently processing API requests.",
		},
	)
	reg.MustRegister(metric.reqInFlight)
	metric.reqTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_total",
			Help:      "Amount of processed API requests.",
		},
		[]string{resourceLabel, methodLabel, codeLabel},
	)
	reg.MustRegister(metric.reqTotal)
	metric.reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_request_duration_seconds",
			Help:      "API request latency distributions.",
		},
		[]string{resourceLabel, methodLabel, failedLabel},
	)
	reg.MustRegister(metric.reqDuration)

	ss, err := loads.Analyzed(restapi.FlatSwaggerJSON, "")
	if err != nil {
		panic(err)
	}
	for method, resources := range ss.Analyzer.Operations() {
		for resource, op := range resources {
			codes := append([]int{}, codeLabels...)
			for code := range op.Responses.StatusCodeResponses {
				codes = append(codes, code)
			}
			for _, code := range codes {
				l := prometheus.Labels{
					resourceLabel: resource,
					methodLabel:   method,
					codeLabel:     strconv.Itoa(code),
				}
				metric.reqTotal.With(l)
			}
			for _, failed := range []string{"true", "false"} {
				l := prometheus.Labels{
					resourceLabel: resource,
					methodLabel:   method,
					failedLabel:   failed,
				}
				metric.reqDuration.With(l)
			}
		}
	}
}
