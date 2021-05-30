// Package config provides configurations for subcommands.
//
// Default values can be obtained from various sources (constants,
// environment variables, etc.) and then overridden by flags.
//
// As configuration is global you can get it only once for safety:
// you can call only one of Get… functions and call it just once.
package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/powerman/appcfg"
	"github.com/spf13/pflag"

	"github.com/powerman/go-service-example/pkg/cobrax"
	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/go-service-example/pkg/netx"
)

// EnvPrefix defines common prefix for environment variables.
const envPrefix = "EXAMPLE_"

// All configurable values of the microservice.
//
// If microservice may runs in different ways (e.g. using CLI subcommands)
// then these subcommands may use subset of these values.
var all = &struct { //nolint:gochecknoglobals // Config is global anyway.
	APIKeyAdmin     appcfg.NotEmptyString `env:"APIKEY_ADMIN"`
	AddrHost        appcfg.NotEmptyString `env:"ADDR_HOST"`
	AddrPort        appcfg.Port           `env:"ADDR_PORT"`
	MetricsAddrPort appcfg.Port           `env:"METRICS_ADDR_PORT"`
	MySQLAddrHost   appcfg.NotEmptyString `env:"MYSQL_ADDR_HOST"`
	MySQLAddrPort   appcfg.Port           `env:"MYSQL_ADDR_PORT"`
	MySQLAuthLogin  appcfg.NotEmptyString `env:"MYSQL_AUTH_LOGIN"`
	MySQLAuthPass   appcfg.String         `env:"MYSQL_AUTH_PASS"`
	MySQLDBName     appcfg.NotEmptyString `env:"MYSQL_DB"`
	GooseMySQLDir   appcfg.NotEmptyString
}{ // Defaults, if any:
	AddrHost:        appcfg.MustNotEmptyString(def.Hostname),
	AddrPort:        appcfg.MustPort("8000"),
	MetricsAddrPort: appcfg.MustPort("9000"),
	MySQLAddrPort:   appcfg.MustPort("3306"),
	MySQLAuthLogin:  appcfg.MustNotEmptyString(def.ProgName),
	MySQLDBName:     appcfg.MustNotEmptyString(def.ProgName),
	GooseMySQLDir:   appcfg.MustNotEmptyString("internal/migrations/mysql"),
}

// FlagSets for all CLI subcommands which use flags to set config values.
type FlagSets struct {
	Serve      *pflag.FlagSet
	GooseMySQL *pflag.FlagSet
}

var fs FlagSets //nolint:gochecknoglobals // Flags are global anyway.

// Init updates config defaults (from env) and setup subcommands flags.
//
// Init must be called once before using this package.
func Init(flagsets FlagSets) error {
	fs = flagsets

	fromEnv := appcfg.NewFromEnv(envPrefix)
	err := appcfg.ProvideStruct(all, fromEnv)
	if err != nil {
		return err
	}

	appcfg.AddPFlag(fs.GooseMySQL, &all.MySQLAddrHost, "mysql.host", "host to connect to MySQL")
	appcfg.AddPFlag(fs.GooseMySQL, &all.MySQLAddrPort, "mysql.port", "port to connect to MySQL")
	appcfg.AddPFlag(fs.GooseMySQL, &all.MySQLDBName, "mysql.dbname", "MySQL database name")
	appcfg.AddPFlag(fs.GooseMySQL, &all.MySQLAuthLogin, "mysql.user", "MySQL username")
	appcfg.AddPFlag(fs.GooseMySQL, &all.MySQLAuthPass, "mysql.pass", "MySQL password")

	appcfg.AddPFlag(fs.Serve, &all.MySQLAddrHost, "mysql.host", "host to connect to MySQL")
	appcfg.AddPFlag(fs.Serve, &all.MySQLAddrPort, "mysql.port", "port to connect to MySQL")
	appcfg.AddPFlag(fs.Serve, &all.MySQLDBName, "mysql.dbname", "MySQL database name")
	appcfg.AddPFlag(fs.Serve, &all.MySQLAuthLogin, "mysql.user", "MySQL username")
	appcfg.AddPFlag(fs.Serve, &all.MySQLAuthPass, "mysql.pass", "MySQL password")
	appcfg.AddPFlag(fs.Serve, &all.AddrHost, "host", "host to serve OpenAPI")
	appcfg.AddPFlag(fs.Serve, &all.AddrPort, "port", "port to serve OpenAPI")
	appcfg.AddPFlag(fs.Serve, &all.MetricsAddrPort, "metrics.port", "port to serve Prometheus metrics")

	return nil
}

// ServeConfig contains configuration for subcommand.
type ServeConfig struct {
	MySQL           *mysql.Config
	GooseMySQLDir   string
	BindAddr        netx.Addr
	BindMetricsAddr netx.Addr
	APIKeyAdmin     string
}

// GetServe validates and returns configuration for subcommand.
func GetServe() (c *ServeConfig, err error) {
	defer cleanup()

	c = &ServeConfig{
		MySQL: def.NewMySQLConfig(def.MySQLConfig{
			Addr:   netx.NewAddr(all.MySQLAddrHost.Value(&err), all.MySQLAddrPort.Value(&err)),
			DBName: all.MySQLDBName.Value(&err),
			User:   all.MySQLAuthLogin.Value(&err),
			Pass:   all.MySQLAuthPass.Value(&err),
		}),
		GooseMySQLDir:   all.GooseMySQLDir.Value(&err),
		BindAddr:        netx.NewAddr(all.AddrHost.Value(&err), all.AddrPort.Value(&err)),
		BindMetricsAddr: netx.NewAddr(all.AddrHost.Value(&err), all.MetricsAddrPort.Value(&err)),
		APIKeyAdmin:     all.APIKeyAdmin.Value(&err),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.Serve, all)
	}
	return c, nil
}

func GetGooseMySQL() (c *cobrax.GooseMySQLConfig, err error) {
	defer cleanup()

	c = &cobrax.GooseMySQLConfig{
		MySQL: def.NewMySQLConfig(def.MySQLConfig{
			Addr:   netx.NewAddr(all.MySQLAddrHost.Value(&err), all.MySQLAddrPort.Value(&err)),
			DBName: all.MySQLDBName.Value(&err),
			User:   all.MySQLAuthLogin.Value(&err),
			Pass:   all.MySQLAuthPass.Value(&err),
		}),
		GooseMySQLDir: all.GooseMySQLDir.Value(&err),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.GooseMySQL, all)
	}
	return c, nil
}

// Cleanup must be called by all Get* functions to ensure second call to
// any of them will panic.
func cleanup() {
	all = nil
}
