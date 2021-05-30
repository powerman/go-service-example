package main

import (
	"context"
	"regexp"

	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"

	"github.com/powerman/go-service-example/api/openapi/restapi"
	"github.com/powerman/go-service-example/internal/app"
	"github.com/powerman/go-service-example/internal/config"
	dal "github.com/powerman/go-service-example/internal/dal/mysql"
	migrations_mysql "github.com/powerman/go-service-example/internal/migrations/mysql"
	"github.com/powerman/go-service-example/internal/srv/openapi"
	"github.com/powerman/go-service-example/pkg/cobrax"
	"github.com/powerman/go-service-example/pkg/concurrent"
	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/go-service-example/pkg/serve"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

var reg = prometheus.NewPedanticRegistry() //nolint:gochecknoglobals // Metrics are global anyway.

type Service struct {
	cfg  *config.ServeConfig
	repo *dal.Repo
	appl *app.App
	srv  *restapi.Server
}

func (s *Service) Init(cmd, serveCmd *cobra.Command) error {
	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(def.ProgName, "_")
	initMetrics(reg, namespace)
	dal.InitMetrics(reg, namespace)
	app.InitMetrics(reg)
	openapi.InitMetrics(reg, namespace)

	gooseMySQLCmd := cobrax.NewGooseMySQLCmd(context.Background(), migrations_mysql.Goose(), config.GetGooseMySQL)
	cmd.AddCommand(gooseMySQLCmd)

	return config.Init(config.FlagSets{
		Serve:      serveCmd.Flags(),
		GooseMySQL: gooseMySQLCmd.Flags(),
	})
}

func (s *Service) RunServe(ctxStartup, ctxShutdown Ctx, shutdown func()) (err error) {
	log := structlog.FromContext(ctxShutdown, nil)
	if s.cfg == nil {
		s.cfg, err = config.GetServe()
	}
	if err != nil {
		return log.Err("failed to get config", "err", err)
	}

	err = concurrent.Setup(ctxStartup, map[interface{}]concurrent.SetupFunc{
		&s.repo: s.connectRepo,
	})
	if err != nil {
		return log.Err("failed to connect", "err", err)
	}

	if s.appl == nil {
		s.appl = app.New(s.repo, app.Config{})
	}
	s.srv, err = openapi.NewServer(s.appl, openapi.Config{
		APIKeyAdmin: s.cfg.APIKeyAdmin,
		Addr:        s.cfg.BindAddr,
	})
	if err != nil {
		return log.Err("failed to openapi.NewServer", "err", err)
	}

	err = concurrent.Serve(ctxShutdown, shutdown,
		s.serveMetrics,
		s.serveOpenAPI,
	)
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	return nil
}

func (s *Service) connectRepo(ctx Ctx) (interface{}, error) {
	return dal.New(ctx, s.cfg.GooseMySQLDir, s.cfg.MySQL)
}

func (s *Service) serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, s.cfg.BindMetricsAddr, reg)
}

func (s *Service) serveOpenAPI(ctx Ctx) error {
	return serve.OpenAPI(ctx, s.srv, "OpenAPI")
}
