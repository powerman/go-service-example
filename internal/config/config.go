// Package config provides configurations for subcommands.
//
// Default values can be obtained from various sources (constants,
// environment variables, etc.) and then overridden by flags.
//
// As configuration is global you can get it only once for safety:
// you can call only one of Get… functions and call it just once.
package config

import (
	"strconv"

	"github.com/powerman/appcfg"
	"github.com/powerman/go-service-goswagger-clean-example/internal/def"
	"github.com/powerman/go-service-goswagger-clean-example/internal/pkg/netx"
	"github.com/spf13/pflag"
)

// envPrefix defines common prefix for environment variables.
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
}{ // Defaults, if any:
	AddrHost:        appcfg.MustNotEmptyString(def.Hostname),
	AddrPort:        appcfg.MustPort(strconv.Itoa(def.Port)),
	MetricsAddrPort: appcfg.MustPort(strconv.Itoa(def.MetricsPort)),
}

// FlagSets for all CLI subcommands which use flags to set config values.
type FlagSets struct {
	Serve *pflag.FlagSet
}

var fs FlagSets //nolint:gochecknoglobals // Flags are global anyway.

// Init updates config defaults (from env) and setup subcommands flags.
func Init(flagsets FlagSets) error {
	fs = flagsets

	fromEnv := appcfg.NewFromEnv(envPrefix)
	err := appcfg.ProvideStruct(all, fromEnv)
	if err != nil {
		return err
	}

	appcfg.AddPFlag(fs.Serve, &all.AddrHost, "host", "host to serve OpenAPI")
	appcfg.AddPFlag(fs.Serve, &all.AddrPort, "port", "port to serve OpenAPI")
	appcfg.AddPFlag(fs.Serve, &all.MetricsAddrPort, "metrics.port", "port to serve Prometheus metrics")

	return nil
}

// ServeConfig contains configuration for subcommand.
type ServeConfig struct {
	APIKeyAdmin string
	Addr        netx.Addr
	MetricsAddr netx.Addr
}

// GetServe validates and returns configuration for subcommand.
func GetServe() (c *ServeConfig, err error) {
	defer cleanup()

	c = &ServeConfig{
		APIKeyAdmin: all.APIKeyAdmin.Value(&err),
		Addr:        netx.NewAddr(all.AddrHost.Value(&err), all.AddrPort.Value(&err)),
		MetricsAddr: netx.NewAddr(all.AddrHost.Value(&err), all.MetricsAddrPort.Value(&err)),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.Serve, all)
	}
	return c, nil
}

// Cleanup must be called by all Get* functions to ensure second call to
// any of them will panic.
func cleanup() {
	all = nil
}
