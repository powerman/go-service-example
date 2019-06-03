package main

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(namespace string) {
	promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "build_info",
			Help:      "A metric with a constant '1' value labeled by build-time details.",
		},
		[]string{"version", "branch", "revision", "date", "build_date", "goversion"},
	).With(prometheus.Labels{
		"version":    gitVersion,
		"branch":     gitBranch,
		"revision":   gitRevision,
		"date":       gitDate,
		"build_date": buildDate,
		"goversion":  runtime.Version(),
	}).Set(1)
}
