// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/crxfoz/seo_metrick_parser/server/api/restapi/operations"

	models "github.com/crxfoz/seo_metrick_parser/server/api/models"
)

//go:generate swagger generate server --target ../api --name SEOParser --spec ../swagger/swagger.yml --principal models.AuthToken --exclude-main

func configureFlags(api *operations.SEOParserAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SEOParserAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "x-token" header is set
	api.APIKeyAuth = func(token string) (*models.AuthToken, error) {
		return nil, errors.NotImplemented("api key auth (ApiKey) x-token from header param [x-token] has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.AddTaskHandler = operations.AddTaskHandlerFunc(func(params operations.AddTaskParams) middleware.Responder {
		return middleware.NotImplemented("operation .AddTask has not yet been implemented")
	})
	api.GetParsersHandler = operations.GetParsersHandlerFunc(func(params operations.GetParsersParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetParsers has not yet been implemented")
	})
	api.GetTaskDataHandler = operations.GetTaskDataHandlerFunc(func(params operations.GetTaskDataParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetTaskData has not yet been implemented")
	})
	api.GetTaskStatusHandler = operations.GetTaskStatusHandlerFunc(func(params operations.GetTaskStatusParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetTaskStatus has not yet been implemented")
	})
	api.GetTasksHandler = operations.GetTasksHandlerFunc(func(params operations.GetTasksParams, principal *models.AuthToken) middleware.Responder {
		return middleware.NotImplemented("operation .GetTasks has not yet been implemented")
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
