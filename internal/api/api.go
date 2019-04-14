package api

import (
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/powerman/go-service-goswagger-clean-example/internal/api/restapi"
	"github.com/powerman/go-service-goswagger-clean-example/internal/api/restapi/op"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/structlog"
	"github.com/rs/cors"
)

type service struct {
	log *structlog.Logger
	app app.App
}

// Config contains configuration for internal API service.
type Config struct {
	Host string
	Port int
}

// Serve listens on the TCP network address cfg.Host:cfg.Port and
// handle requests on incoming connections.
func Serve(log *structlog.Logger, application app.App, cfg Config) error {
	svc := &service{
		log: log,
		app: application,
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return err
	}

	api := op.NewAddressBookAPI(swaggerSpec)
	api.Logger = log.Printf
	api.APIKeyAuth = svc.authenticate
	api.APIAuthorizer = runtime.AuthorizerFunc(svc.authorize)
	api.GetContactsHandler = op.GetContactsHandlerFunc(svc.getContacts)
	api.PostContactsHandler = op.PostContactsHandlerFunc(svc.postContacts)

	server := restapi.NewServer(api)
	defer log.WarnIfFail(server.Shutdown)

	server.Host = cfg.Host
	server.Port = cfg.Port

	server.SetHandler(setupGlobalMiddlewares(api.Serve(setupMiddlewares)))

	log.Info("protocol", "version", swaggerSpec.Spec().Info.Version)
	return server.Serve()
}

// The middleware configuration happens before anything.
// This middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddlewares(handler http.Handler) http.Handler {
	handleCORS := cors.AllowAll().Handler
	return handleCORS(handler)
}

// The middleware configuration is for the handler executors.
// These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}
