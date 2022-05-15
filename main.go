package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
)

const (
	appName      = "my-app"
	hdrRequestId = "X-Request-ID"
	hdrTracingId = "X-Tracing-ID"
	appAddr      = "0.0.0.0:4040"
	logFile      = "_logs/observability.log"
	tracingUrl   = "http://localhost:14268/api/traces"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// Set JSON formatter. Comment out this line to have the text output.
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Set logging to a file. Comment out following 2 lines to log on the console.
	f := getLogFile()
	logrus.SetOutput(f)
	// Register Prometheus metrics
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
	// Set tracing provider
	tp, err := tracerProvider(tracingUrl)
	if err != nil {
		log.Fatal(err)
	}
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
}

func main() {
	log := funcLog("main")
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.Handle("/metrics", promhttp.Handler())
	r.Use(tracingMiddleware)
	r.Use(metricsMiddleware)
	r.Use(loggingMiddleware)
	log.Infof("starting observability app on: %s", appAddr)
	http.ListenAndServe(appAddr, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer(appName).Start(r.Context(), "homeHandler")
	defer span.End()
	log := requestLog("homeHandler", r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := make(map[string]string)
	resp["message"] = "Observability check: 👌"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("can't marshal json: %v", err)
	}
	log.Infof("writing response with status: %d", http.StatusOK)
	w.Write(jsonResp)
	return
}
