package main

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func requestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := requestLog("requestLogMiddleware", getRequestId(r))
		log.Infof("serving request: %s %s%s", r.Method, r.Host, r.RequestURI)
		log.Debugf("user agent: %s", r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

func getRequestId(r *http.Request) string {
	requestId := r.Header.Get(hdrRequestId)
	if requestId == "" {
		requestId = uuid.New().String()
		r.Header.Set(hdrRequestId, requestId)
		log := requestLog("getRequestId", requestId)
		log.Warnf("header %s is empty, no request id has been provided", hdrRequestId)
	}
	return requestId
}

func requestLog(f, id string) *logrus.Entry {
	return funcLog(f).WithField("requestId", id)
}

func funcLog(f string) *logrus.Entry {
	return logrus.WithField("func", f)
}

func getLogFile() *os.File {
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		funcLog("getLogFile").Fatalf("log file %s can't be created: %v", logFile, err)
	}
	return f
}
