package app

import (
	"github.com/powerman/go-service-example/pkg/def"
	"github.com/prometheus/client_golang/prometheus"
)

//nolint:gochecknoglobals // Metrics are global anyway.
var (
	Metric def.Metrics // Common metrics used by all packages.
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(reg *prometheus.Registry) {
	Metric = def.NewMetrics(reg)
}
