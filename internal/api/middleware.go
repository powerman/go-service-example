package api

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/powerman/go-service-goswagger-clean-example/internal/def"
	"github.com/powerman/structlog"
	"github.com/rs/cors"
	"github.com/sebest/xff"
)

type middlewareFunc func(http.Handler) http.Handler

// Provide a logger configured using request's context.
//
// Usually it should be first middleware.
func makeLogger(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			remote := xff.GetRemoteAddr(r)
			log := structlog.New(
				def.LogRemote, remote,
				def.LogHTTPStatus, "",
				def.LogHTTPMethod, r.Method,
				def.LogFunc, strings.TrimPrefix(r.URL.Path, basePath),
			)
			r = r.WithContext(structlog.NewContext(r.Context(), log))

			next.ServeHTTP(w, r)
		})
	}
}

// go-swagger responders panic on error while writing response to client,
// this shouldn't result in crash - unlike a real, reasonable panic.
//
// Usually it should be second middleware (after logger).
func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			switch err := recover(); err := err.(type) {
			default:
				log := structlog.FromContext(r.Context(), nil)
				log.PrintErr(err, structlog.KeyStack, structlog.Auto)
				os.Exit(2)
			case nil:
			case net.Error:
				log := structlog.FromContext(r.Context(), nil)
				log.PrintErr(err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func handleCORS(next http.Handler) http.Handler {
	return cors.AllowAll().Handler(next)
}

func accesslog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := wrapResponseWriter(w)

		next.ServeHTTP(ww, r)

		log := structlog.FromContext(r.Context(), nil)
		if code := ww.StatusCode(); code < 500 {
			log.Info("request handled", def.LogHTTPStatus, code)
		} else {
			log.PrintErr("failed to handle request", def.LogHTTPStatus, code)
		}
	})
}
