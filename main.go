package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	hdrRequestId = "X-Request-ID"
	appAddr      = "0.0.0.0:4040"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	log := funcLog("main")
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.Use(requestIdMiddleware)
	r.Use(requestLogMiddleware)
	log.Infof("starting observability app on: %s", appAddr)
	http.ListenAndServe(appAddr, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log := funcLog("homeHandler")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := make(map[string]string)
	resp["message"] = "Observability check: 👌"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("can't marshal json: %w", err)
	}
	log.Infof("writing response with status: %d", http.StatusOK)
	w.Write(jsonResp)
	return
}

func requestIdMiddleware(next http.Handler) http.Handler {
	log := funcLog("requestIdMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get(hdrRequestId)
		if requestId == "" {
			log.Warnf("header %s is empty, no request id has been provided", hdrRequestId)
		} else {
			log.Infof("client request id: %s", requestId)
		}
		next.ServeHTTP(w, r)
	})
}

func requestLogMiddleware(next http.Handler) http.Handler {
	log := funcLog("requestLogMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("serving request: %s %s%s", r.Method, r.Host, r.RequestURI)
		log.Debugf("user agent: %s", r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

func logRequest(r *http.Request) {
}

func funcLog(f string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"func": f,
	})
}
