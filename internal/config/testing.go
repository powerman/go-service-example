package config

import (
	"github.com/powerman/go-service-goswagger-clean-example/internal/pkg/netx"
	"github.com/powerman/must"
	"github.com/spf13/pflag"
)

// MustGetServeTest returns config suitable for use in tests.
func MustGetServeTest() *ServeConfig {
	err := Init(FlagSets{
		Serve: pflag.NewFlagSet("", pflag.ContinueOnError),
	})
	must.NoErr(err)
	cfg, err := GetServe()
	must.NoErr(err)

	const host = "localhost"
	cfg.Addr = netx.NewAddr(host, netx.UnusedTCPPort(host))
	cfg.MetricsAddr = netx.NewAddr(host, 0)

	return cfg
}
