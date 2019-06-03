package api

import (
	"context"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/powerman/go-service-goswagger-clean-example/internal/api/restapi"
	"github.com/powerman/go-service-goswagger-clean-example/internal/api/restapi/op"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/go-service-goswagger-clean-example/internal/def"
	"github.com/powerman/structlog"
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

	// The middleware executes before anything.
	globalMiddlewares := func(handler http.Handler) http.Handler {
		logger := makeLogger(swaggerSpec.BasePath())
		return logger(recovery(handleCORS(handler)))
	}
	// The middleware executes after serving /swagger.json and routing,
	// but before authentication, binding and validation.
	middlewares := func(handler http.Handler) http.Handler {
		accesslog := makeAccessLog(swaggerSpec.BasePath())
		return accesslog(handler)
	}
	server.SetHandler(globalMiddlewares(api.Serve(middlewares)))

	log.Info("protocol", "version", swaggerSpec.Spec().Info.Version)
	return server.Serve()
}

func fromRequest(r *http.Request, auth *app.Auth) (context.Context, *structlog.Logger) {
	ctx := r.Context()
	log := structlog.FromContext(ctx, nil).SetDefaultKeyvals(def.LogUser, auth.UserID)
	return ctx, log
}
