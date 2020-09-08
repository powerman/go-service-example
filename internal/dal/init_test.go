package dal_test

import (
	"context"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/go-service-goswagger-clean-example/pkg/def"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	def.Init()
	reg := prometheus.NewPedanticRegistry()
	app.InitMetrics(reg)
	check.TestMain(m)
}

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	ctx = context.Background()
)
