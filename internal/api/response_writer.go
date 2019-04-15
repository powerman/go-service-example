package api

import "net/http"

type responseWriter interface {
	http.ResponseWriter
	StatusCode() int
}

type responseWriter0 struct { // Compatible with net/http.ResponseWriter for HTTP 1.x and 2.
	statusCode int
	http.ResponseWriter
	http.CloseNotifier
	http.Flusher
}

type responseWriter1 struct { // Compatible with net/http.ResponseWriter for HTTP 1.x.
	responseWriter0
	http.Hijacker
}

type responseWriter2 struct { // Compatible with net/http.ResponseWriter for HTTP 2.
	responseWriter0
	http.Pusher
}

func wrapResponseWriter(w http.ResponseWriter) responseWriter {
	rw0 := responseWriter0{
		statusCode:     http.StatusOK,
		ResponseWriter: w,
		CloseNotifier:  w.(http.CloseNotifier), //nolint:staticcheck
		Flusher:        w.(http.Flusher),
	}
	switch w := w.(type) {
	case http.Hijacker:
		return &responseWriter1{
			responseWriter0: rw0,
			Hijacker:        w,
		}
	case http.Pusher:
		return &responseWriter2{
			responseWriter0: rw0,
			Pusher:          w,
		}
	default:
		panic("unknown implementation of http.ResponseWriter")
	}
}

func (rw *responseWriter0) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter0) StatusCode() int {
	return rw.statusCode
}
