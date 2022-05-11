package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of http requests.",
	},
	[]string{"method", "path", "app"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_response_status",
		Help: "Status of http response",
	},
	[]string{"status", "app"},
)

var httpDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of http requests.",
	},
	[]string{"method", "path", "app"},
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(
			r.Method,
			r.RequestURI,
			appName,
		))
		defer timer.ObserveDuration()
		rw := newResponseWriter(w)
		next.ServeHTTP(w, r)
		totalRequests.WithLabelValues(
			r.Method,
			r.RequestURI,
			appName,
		).Inc()
		responseStatus.WithLabelValues(
			strconv.Itoa(rw.statusCode),
			appName,
		).Inc()
	})
}
