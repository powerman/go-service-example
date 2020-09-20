package dal

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-service-example/internal/app"
	"github.com/powerman/go-service-example/pkg/repo"
)

var metric repo.Metrics //nolint:gochecknoglobals // Metrics are global anyway.

func InitMetrics(reg *prometheus.Registry, namespace string) {
	const subsystem = "dal_mysql"

	metric = repo.NewMetrics(reg, namespace, subsystem, new(app.Repo))
}
