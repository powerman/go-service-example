package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/powerman/must"
	"github.com/spf13/pflag"

	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/go-service-example/pkg/netx"
)

// MustGetServeTest returns config suitable for use in tests.
func MustGetServeTest() *ServeConfig {
	err := Init(FlagSets{
		Serve:      pflag.NewFlagSet("", pflag.ContinueOnError),
		GooseMySQL: pflag.NewFlagSet("", pflag.ContinueOnError),
	})
	must.NoErr(err)
	cfg, err := GetServe()
	must.NoErr(err)

	cfg.MySQL.Timeout = def.TestTimeout

	const host = "localhost"
	cfg.BindAddr = netx.NewAddr(host, netx.UnusedTCPPort(host))
	cfg.BindMetricsAddr = netx.NewAddr(host, 0)

	rootDir, err := os.Getwd()
	must.NoErr(err)
	for _, err := os.Stat(filepath.Join(rootDir, "go.mod")); os.IsNotExist(err) && filepath.Dir(rootDir) != rootDir; _, err = os.Stat(filepath.Join(rootDir, "go.mod")) {
		rootDir = filepath.Dir(rootDir)
	}

	for _, path := range []*string{
		&cfg.GooseMySQLDir,
	} {
		if !strings.HasPrefix(*path, "/") {
			*path = filepath.Join(rootDir, *path)
		}
	}

	return cfg
}
