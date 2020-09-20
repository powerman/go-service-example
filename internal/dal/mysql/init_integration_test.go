// +build integration

package dal_test

import (
	"context"
	"runtime"
	"strings"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/go-service-example/internal/app"
	"github.com/powerman/go-service-example/internal/config"
	dal "github.com/powerman/go-service-example/internal/dal/mysql"
	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/mysqlx"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMain(m *testing.M) {
	def.Init()
	reg := prometheus.NewPedanticRegistry()
	app.InitMetrics(reg)
	dal.InitMetrics(reg, "test")
	cfg = config.MustGetServeTest()
	check.TestMain(m)
}

var (
	ctx = context.Background()
	cfg *config.ServeConfig
)

type tLogger check.C

func (t tLogger) Print(args ...interface{}) { t.Log(args...) }

func newTestRepo(t *check.C) (cleanup func(), r *dal.Repo) {
	t.Helper()

	pc, _, _, _ := runtime.Caller(1)
	suffix := runtime.FuncForPC(pc).Name()
	suffix = suffix[:strings.LastIndex(suffix, ".")]
	suffix += "_" + t.Name()

	tempDBCfg, cleanupDB, err := mysqlx.EnsureTempDB(tLogger(*t), suffix, cfg.MySQL)
	t.Must(t.Nil(err))
	r, err = dal.New(ctx, cfg.MySQLGooseDir, tempDBCfg)
	t.Must(t.Nil(err))

	cleanup = func() {
		r.Close()
		cleanupDB()
	}
	return cleanup, r
}
