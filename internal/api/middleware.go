package api

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/powerman/go-service-goswagger-clean-example/internal/def"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
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
				def.Metric.PanicsTotal.Inc()
				log := structlog.FromContext(r.Context(), nil)
				log.PrintErr(err, structlog.KeyStack, structlog.Auto)
				w.WriteHeader(http.StatusInternalServerError)
			case nil:
			case net.Error:
				log := structlog.FromContext(r.Context(), nil)
				log.PrintErr(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func handleCORS(next http.Handler) http.Handler {
	return cors.AllowAll().Handler(next)
}

func makeAccessLog(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			metric.reqInFlight.Inc()
			defer metric.reqInFlight.Dec()

			start := time.Now()
			ww := wrapResponseWriter(w)

			next.ServeHTTP(ww, r)

			code := ww.StatusCode()

			l := prometheus.Labels{
				resourceLabel: strings.TrimPrefix(r.URL.Path, basePath),
				methodLabel:   r.Method,
				codeLabel:     strconv.Itoa(code),
			}
			metric.reqTotal.With(l).Inc()
			metric.reqDuration.With(l).Observe(time.Since(start).Seconds())

			log := structlog.FromContext(r.Context(), nil)
			if code < 500 {
				log.Info("handled", def.LogHTTPStatus, code)
			} else {
				log.PrintErr("failed to handle", def.LogHTTPStatus, code)
			}
		})
	}
}
