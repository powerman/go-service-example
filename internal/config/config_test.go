package config

import (
	"os"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/go-service-example/pkg/netx"
)

func Test(t *testing.T) {
	want := &ServeConfig{
		MySQL: def.NewMySQLConfig(def.MySQLConfig{
			Addr:   netx.NewAddr("localhost", 3306),
			DBName: "config",
			User:   "config",
			Pass:   "",
		}),
		GooseMySQLDir:   "internal/migrations/mysql",
		BindAddr:        netx.NewAddr(def.Hostname, 8000),
		BindMetricsAddr: netx.NewAddr(def.Hostname, 9000),
		APIKeyAdmin:     "admin",
	}

	t.Run("required", func(tt *testing.T) {
		t := check.T(tt)
		require(t, "APIKeyAdmin")
		os.Setenv("EXAMPLE_APIKEY_ADMIN", "admin")
		require(t, "MySQLAuthPass")
		os.Setenv("EXAMPLE_MYSQL_AUTH_PASS", "")
		require(t, "MySQLAddrHost")
		os.Setenv("EXAMPLE_MYSQL_ADDR_HOST", "localhost")
	})
	t.Run("default", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe()
		t.Nil(err)
		t.DeepEqual(c, want)
	})
	t.Run("constraint", func(tt *testing.T) {
		t := check.T(tt)
		constraint(t, "EXAMPLE_ADDR_PORT", "x", `^AddrPort .* invalid syntax`)
		constraint(t, "EXAMPLE_METRICS_ADDR_PORT", "x", `^MetricsAddrPort .* invalid syntax`)
		constraint(t, "EXAMPLE_MYSQL_ADDR_HOST", "", `^MySQLAddrHost .* empty`)
		constraint(t, "EXAMPLE_MYSQL_ADDR_PORT", "x", `^MySQLAddrPort .* invalid syntax`)
		constraint(t, "EXAMPLE_MYSQL_AUTH_LOGIN", "", `^MySQLAuthLogin .* empty`)
		constraint(t, "EXAMPLE_MYSQL_DB", "", `^MySQLDBName .* empty`)
	})
	t.Run("env", func(tt *testing.T) {
		t := check.T(tt)
		os.Setenv("EXAMPLE_APIKEY_ADMIN", "admin3")
		os.Setenv("EXAMPLE_ADDR_HOST", "localhost3")
		os.Setenv("EXAMPLE_ADDR_PORT", "8003")
		os.Setenv("EXAMPLE_METRICS_ADDR_PORT", "9003")
		os.Setenv("EXAMPLE_MYSQL_ADDR_HOST", "mysql3")
		os.Setenv("EXAMPLE_MYSQL_ADDR_PORT", "33306")
		os.Setenv("EXAMPLE_MYSQL_AUTH_LOGIN", "user3")
		os.Setenv("EXAMPLE_MYSQL_AUTH_PASS", "pass3")
		os.Setenv("EXAMPLE_MYSQL_DB", "db3")
		c, err := testGetServe()
		t.Nil(err)
		want.MySQL = def.NewMySQLConfig(def.MySQLConfig{
			Addr:   netx.NewAddr("mysql3", 33306),
			DBName: "db3",
			User:   "user3",
			Pass:   "pass3",
		})
		want.BindAddr = netx.NewAddr("localhost3", 8003)
		want.BindMetricsAddr = netx.NewAddr("localhost3", 9003)
		want.APIKeyAdmin = "admin3"
		t.DeepEqual(c, want)
	})
	t.Run("flag", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe(
			"--mysql.host=mysql4",
			"--mysql.port=43306",
			"--mysql.user=user4",
			"--mysql.pass=pass4",
			"--mysql.dbname=db4",
			"--host=localhost4",
			"--port=8004",
			"--metrics.port=9004",
		)
		t.Nil(err)
		want.MySQL = def.NewMySQLConfig(def.MySQLConfig{
			Addr:   netx.NewAddr("mysql4", 43306),
			DBName: "db4",
			User:   "user4",
			Pass:   "pass4",
		})
		want.BindAddr = netx.NewAddr("localhost4", 8004)
		want.BindMetricsAddr = netx.NewAddr("localhost4", 9004)
		t.DeepEqual(c, want)
	})
	t.Run("cleanup", func(tt *testing.T) {
		t := check.T(tt)
		t.Panic(func() { GetServe() })
	})
}
