// Package openapi implements OpenAPI server.
package openapi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/restapi"
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/restapi/op"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/go-service-goswagger-clean-example/pkg/def"
	"github.com/powerman/go-service-goswagger-clean-example/pkg/netx"
	"github.com/powerman/structlog"
	"github.com/sebest/xff"
)

type (
	// Ctx is a synonym for convenience.
	Ctx = context.Context
	// Log is a synonym for convenience.
	Log = *structlog.Logger
	// Config contains configuration for OpenAPI server.
	Config struct {
		APIKeyAdmin string
		Addr        netx.Addr
		BasePath    string
	}
	server struct {
		app app.Appl
		cfg Config
	}
)

// NewServer returns OpenAPI server configured to listen on the TCP network
// address cfg.Host:cfg.Port and handle requests on incoming connections.
func NewServer(appl app.Appl, cfg Config) (*restapi.Server, error) {
	srv := &server{
		app: appl,
		cfg: cfg,
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return nil, fmt.Errorf("load embedded swagger spec: %w", err)
	}
	if cfg.BasePath == "" {
		cfg.BasePath = swaggerSpec.BasePath()
	}
	swaggerSpec.Spec().BasePath = cfg.BasePath

	api := op.NewAddressBookAPI(swaggerSpec)
	api.Logger = structlog.New(structlog.KeyUnit, "swagger").Printf
	api.APIKeyAuth = srv.authenticate
	api.APIAuthorizer = runtime.AuthorizerFunc(srv.authorize)

	api.ListContactsHandler = op.ListContactsHandlerFunc(srv.listContacts)
	api.AddContactHandler = op.AddContactHandlerFunc(srv.addContact)

	server := restapi.NewServer(api)
	server.Host = cfg.Addr.Host()
	server.Port = cfg.Addr.Port()

	// The middleware executes before anything.
	globalMiddlewares := func(handler http.Handler) http.Handler {
		xffmw, _ := xff.Default()
		logger := makeLogger(cfg.BasePath)
		accesslog := makeAccessLog(cfg.BasePath)
		redocOpts := middleware.RedocOpts{
			BasePath: cfg.BasePath,
			SpecURL:  path.Join(cfg.BasePath, "/swagger.json"),
		}
		return noCache(xffmw.Handler(logger(recovery(accesslog(
			middleware.Spec(cfg.BasePath, restapi.FlatSwaggerJSON,
				middleware.Redoc(redocOpts,
					handleCORS(handler))))))))
	}
	// The middleware executes after serving /swagger.json and routing,
	// but before authentication, binding and validation.
	middlewares := func(handler http.Handler) http.Handler {
		return handler
	}
	server.SetHandler(globalMiddlewares(api.Serve(middlewares)))

	log := structlog.New()
	log.Info("OpenAPI protocol", "version", swaggerSpec.Spec().Info.Version)
	return server, nil
}

func fromRequest(r *http.Request, auth *app.Auth) (Ctx, Log, string) { //nolint:unparam // Some results may be unused yet.
	ctx := r.Context()
	userID := ""
	if auth != nil {
		userID = auth.UserID
	}
	log := structlog.FromContext(ctx, nil).SetDefaultKeyvals(def.LogUserID, userID)
	remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ctx, log, remoteIP
}
