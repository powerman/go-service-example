// Example microservice.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/powerman/appcfg"
	"github.com/powerman/structlog"
	"github.com/spf13/cobra"

	"github.com/powerman/go-service-example/pkg/cobrax"
	"github.com/powerman/go-service-example/pkg/def"
)

//nolint:gochecknoglobals // Main.
var (
	svc = &Service{}

	log = structlog.New(structlog.KeyUnit, "main")

	logLevel = appcfg.MustOneOfString("debug", []string{"debug", "info", "warn", "err"})
	rootCmd  = &cobra.Command{
		Use:           def.ProgName,
		Short:         "Example microservice",
		Version:       fmt.Sprintf("%s %s", def.Version(), runtime.Version()),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          cobrax.RequireFlagOrCommand,
	}

	serveStartupTimeout  = appcfg.MustDuration("3s") // must be less than swarm's deploy.update_config.monitor
	serveShutdownTimeout = appcfg.MustDuration("9s") // `docker stop` use 10s between SIGTERM and SIGKILL
	serveCmd             = &cobra.Command{
		Use:   "serve",
		Short: "Starts microservice",
		Args:  cobra.NoArgs,
		RunE:  runServeWithGracefulShutdown,
	}
)

func main() {
	_ = os.Unsetenv("GO_TEST_DISABLE_SENSITIVE")
	err := def.Init()
	if err != nil {
		log.Fatalf("failed to get defaults: %s", err)
	}

	err = svc.Init(rootCmd, serveCmd)
	if err != nil {
		log.Fatalf("failed to init service: %s", err)
	}

	rootCmd.PersistentFlags().Var(&logLevel, "log.level", "log level [debug|info|warn|err]")
	serveCmd.Flags().Var(&serveStartupTimeout, "timeout.startup", "must be less than swarm's deploy.update_config.monitor")
	serveCmd.Flags().Var(&serveShutdownTimeout, "timeout.shutdown", "must be less than 10s used by 'docker stop' between SIGTERM and SIGKILL")
	rootCmd.AddCommand(serveCmd)

	cobra.OnInitialize(func() {
		structlog.DefaultLogger.SetLogLevel(structlog.ParseLevel(logLevel.String()))
	})
	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func runServeWithGracefulShutdown(_ *cobra.Command, _ []string) error {
	log.Info("started", "version", def.Version())
	defer log.Info("finished", "version", def.Version())

	ctxStartup, cancel := context.WithTimeout(context.Background(), serveStartupTimeout.Value(nil))
	defer cancel()

	ctxShutdown, shutdown := context.WithCancel(context.Background())
	ctxShutdown, _ = signal.NotifyContext(ctxShutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() {
		<-ctxShutdown.Done()
		time.Sleep(serveShutdownTimeout.Value(nil))
		log.PrintErr("failed to graceful shutdown", "version", def.Version())
		os.Exit(1)
	}()

	return svc.RunServe(ctxStartup, ctxShutdown, shutdown)
}
