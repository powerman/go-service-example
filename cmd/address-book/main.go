// Example swagger service.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/powerman/go-service-goswagger-clean-example/internal/api"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/go-service-goswagger-clean-example/internal/def"
	"github.com/powerman/go-service-goswagger-clean-example/internal/flags"
	"github.com/powerman/structlog"
)

//nolint:gochecknoglobals
var (
	cmd = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	ver string // set by ./build
	log = structlog.New()
	cfg struct {
		version  bool
		logLevel string
		api      api.Config
	}
)

// Init provides common initialization for both app and tests.
func Init() {
	def.Init()

	flag.BoolVar(&cfg.version, "version", false, "print version")
	flag.StringVar(&cfg.logLevel, "log.level", "debug", "log `level` (debug|info|warn|err)")
	flag.StringVar(&cfg.api.Host, "host", def.Host, "listen on `host`")
	flag.IntVar(&cfg.api.Port, "port", def.Port, "listen on `port` (>0)")

	log.SetDefaultKeyvals(
		structlog.KeyUnit, "main",
	)
}

func main() {
	Init()
	flag.Parse()

	switch {
	case cfg.api.Host == "":
		flags.FatalFlagValue("required", "host", cfg.api.Host)
	case cfg.api.Port <= 0: // Free nginx doesn't support dynamic ports.
		flags.FatalFlagValue("must be > 0", "port", cfg.api.Port)
	case cfg.version: // Must be checked after all other flags for ease testing.
		fmt.Println(cmd, ver, runtime.Version())
		os.Exit(0)
	}

	// Wrong log.level is not fatal, it will be reported and set to "debug".
	structlog.DefaultLogger.SetLogLevel(structlog.ParseLevel(cfg.logLevel))
	log.Info("started", "version", ver)

	a := app.New()
	err := api.Serve(log, a, cfg.api)
	if err != nil {
		log.Fatal(err)
	}
}
