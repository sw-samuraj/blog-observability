package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	appName      = "my-app"
	hdrRequestId = "X-Request-ID"
	appAddr      = "0.0.0.0:4040"
	logFile      = "_logs/observability.log"
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
}

func main() {
	log := funcLog("main")
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.Handle("/metrics", promhttp.Handler())
	r.Use(requestLogMiddleware)
	r.Use(prometheusMiddleware)
	log.Infof("starting observability app on: %s", appAddr)
	http.ListenAndServe(appAddr, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	requestId := getRequestId(r)
	log := requestLog("homeHandler", requestId)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(hdrRequestId, requestId)
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
