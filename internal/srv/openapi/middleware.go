package openapi

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/go-openapi/swag"
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/model"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/go-service-goswagger-clean-example/pkg/def"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
)

type middlewareFunc func(http.Handler) http.Handler

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", "0")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		next.ServeHTTP(w, r)
	})
}

// Provide a logger configured using request's context.
//
// Usually it should be one of the first (but after xff, if used) middleware.
func makeLogger(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := structlog.New(
				def.LogRemote, r.RemoteAddr,
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
		panicked := true
		defer func() {
			if panicked {
				const code = http.StatusInternalServerError
				switch err := recover(); err := err.(type) {
				default:
					app.Metric.PanicsTotal.Inc()
					log := structlog.FromContext(r.Context(), nil)
					log.PrintErr("panic", def.LogHTTPStatus, code, "err", err, structlog.KeyStack, structlog.Auto)
					middlewareError(w, code, "internal error")
				case net.Error:
					log := structlog.FromContext(r.Context(), nil)
					log.PrintErr("recovered", def.LogHTTPStatus, code, "err", err)
					middlewareError(w, code, "internal error")
				}
			}
		}()
		next.ServeHTTP(w, r)
		panicked = false
	})
}

func makeAccessLog(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			metric.reqInFlight.Inc()
			defer metric.reqInFlight.Dec()

			m := httpsnoop.CaptureMetrics(next, w, r)

			l := prometheus.Labels{
				resourceLabel: strings.TrimPrefix(r.URL.Path, basePath),
				methodLabel:   r.Method,
				codeLabel:     strconv.Itoa(m.Code),
			}
			metric.reqTotal.With(l).Inc()
			l = prometheus.Labels{
				resourceLabel: strings.TrimPrefix(r.URL.Path, basePath),
				methodLabel:   r.Method,
				failedLabel:   strconv.FormatBool(m.Code >= http.StatusInternalServerError),
			}
			metric.reqDuration.With(l).Observe(m.Duration.Seconds())

			log := structlog.FromContext(r.Context(), nil)
			if m.Code < http.StatusInternalServerError {
				log.Info("handled", def.LogHTTPStatus, m.Code)
			} else {
				log.PrintErr("failed to handle", def.LogHTTPStatus, m.Code)
			}
		})
	}
}

func handleCORS(next http.Handler) http.Handler {
	return cors.AllowAll().Handler(next)
}

// MiddlewareError is not a middleware, it's a helper for returning errors
// from middleware.
func middlewareError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&model.Error{
		Code:    swag.Int32(int32(code)),
		Message: swag.String(msg),
	})
}
