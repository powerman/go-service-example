// +build integration

package migrations_test

import (
	"context"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-service-example/internal/config"
	migrations "github.com/powerman/go-service-example/internal/migrations/mysql"
	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/go-service-example/pkg/migrate"
)

var cfg *config.ServeConfig

func TestMain(m *testing.M) {
	def.Init()
	cfg = config.MustGetServeTest()
	check.TestMain(m)
}

func Test(tt *testing.T) {
	t := check.T(tt)
	ctx := context.Background()
	migrate.MySQLUpDownTest(t, ctx, migrations.Goose(), ".", cfg.MySQL)
}
