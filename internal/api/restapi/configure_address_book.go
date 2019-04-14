// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"

	"github.com/powerman/go-service-goswagger-clean-example/internal/api/restapi/op"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
)

//go:generate swagger generate server --target ../../api --name AddressBook --spec ../../../swagger.yml --api-package op --model-package model --principal app.Auth --exclude-main --strict

func configureFlags(api *op.AddressBookAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *op.AddressBookAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "API-Key" header is set
	api.APIKeyAuth = func(token string) (*app.Auth, error) {
		return nil, errors.NotImplemented("api key auth (api_key) API-Key from header param [API-Key] has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.GetContactsHandler = op.GetContactsHandlerFunc(func(params op.GetContactsParams, principal *app.Auth) op.GetContactsResponder {
		return op.GetContactsNotImplemented()
	})
	api.PostContactsHandler = op.PostContactsHandlerFunc(func(params op.PostContactsParams, principal *app.Auth) op.PostContactsResponder {
		return op.PostContactsNotImplemented()
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
