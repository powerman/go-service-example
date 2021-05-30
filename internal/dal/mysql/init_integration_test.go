// +build integration

package dal_test

import (
	"context"
	"runtime"
	"strings"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/mysqlx"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-service-example/internal/app"
	"github.com/powerman/go-service-example/internal/config"
	dal "github.com/powerman/go-service-example/internal/dal/mysql"
	"github.com/powerman/go-service-example/pkg/def"
)

func TestMain(m *testing.M) {
	def.Init()
	reg := prometheus.NewPedanticRegistry()
	app.InitMetrics(reg)
	dal.InitMetrics(reg, "test")
	cfg = config.MustGetServeTest()
	check.TestMain(m)
}

type tLogger check.C

func (t tLogger) Print(args ...interface{}) { t.Log(args...) }

var (
	ctx = context.Background()
	cfg *config.ServeConfig
)

func newTestRepo(t *check.C) *dal.Repo {
	t.Helper()

	pc, _, _, _ := runtime.Caller(1)
	suffix := runtime.FuncForPC(pc).Name()
	suffix = suffix[:strings.LastIndex(suffix, ".")]
	suffix += "_" + t.Name()

	tempDBCfg, cleanupDB, err := mysqlx.EnsureTempDB(tLogger(*t), suffix, cfg.MySQL)
	t.Must(t.Nil(err))
	t.Cleanup(cleanupDB)
	r, err := dal.New(ctx, cfg.GooseMySQLDir, tempDBCfg)
	t.Must(t.Nil(err))
	t.Cleanup(r.Close)

	return r
}
