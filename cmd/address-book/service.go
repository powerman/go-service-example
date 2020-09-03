package main

import (
	"context"
	"regexp"

	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/restapi"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/go-service-goswagger-clean-example/internal/config"
	"github.com/powerman/go-service-goswagger-clean-example/internal/dal"
	"github.com/powerman/go-service-goswagger-clean-example/internal/def"
	"github.com/powerman/go-service-goswagger-clean-example/internal/pkg/concurrent"
	"github.com/powerman/go-service-goswagger-clean-example/internal/pkg/serve"
	"github.com/powerman/go-service-goswagger-clean-example/internal/srv/openapi"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

var reg = prometheus.NewPedanticRegistry() //nolint:gochecknoglobals // Metrics are global anyway.

type service struct {
	cfg  *config.ServeConfig
	repo *dal.Repo
	appl *app.App
	srv  *restapi.Server
}

func initService(_, serveCmd *cobra.Command) error {
	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(def.ProgName, "_")
	initMetrics(reg, namespace)
	app.InitMetrics(reg)
	openapi.InitMetrics(reg, namespace)

	return config.Init(config.FlagSets{
		Serve: serveCmd.Flags(),
	})
}

func (s *service) runServe(ctxStartup, ctxShutdown Ctx, shutdown func()) (err error) {
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
		s.appl = app.New(s.repo)
	}
	s.srv, err = openapi.NewServer(s.appl, openapi.Config{
		APIKeyAdmin: s.cfg.APIKeyAdmin,
		Addr:        s.cfg.Addr,
	})
	if err != nil {
		return err
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

func (s *service) connectRepo(ctx Ctx) (interface{}, error) {
	return dal.New(ctx)
}

func (s *service) serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, s.cfg.MetricsAddr, reg)
}

func (s *service) serveOpenAPI(ctx Ctx) error {
	return serve.OpenAPI(ctx, s.srv, "OpenAPI")
}
