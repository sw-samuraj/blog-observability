package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"net/http"
	"os"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanCtx, span := otel.Tracer(appName).Start(r.Context(), "loggingMiddleware")
		defer span.End()
		log := requestLog("loggingMiddleware", r)
		log.Infof("serving request: %s %s%s", r.Method, r.Host, r.RequestURI)
		log.Debugf("user agent: %s", r.UserAgent())
		spannedRequest := r.WithContext(spanCtx)
		w.Header().Set(hdrRequestId, getRequestId(r))
		next.ServeHTTP(w, spannedRequest)
	})
}

func requestLog(f string, r *http.Request) *logrus.Entry {
	rid := getRequestId(r)
	tid := getTracingId(r)
	return funcLog(f).WithFields(logrus.Fields{
		"requestId": rid,
		"tracingId": tid,
	})
}

func getRequestId(r *http.Request) string {
	requestId := r.Header.Get(hdrRequestId)
	if requestId == "" {
		requestId = uuid.New().String()
		r.Header.Set(hdrRequestId, requestId)
		log := requestLog("getRequestId", r)
		log.Warnf("header %s is empty, no request id has been provided", hdrRequestId)
	}
	return requestId
}

func getTracingId(r *http.Request) string {
	return r.Header.Get(hdrTracingId)
}

func funcLog(f string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":  appName,
		"func": f,
	})
}

func getLogFile() *os.File {
	file := logFile
	if appName != defaultAppName {
		file = fmt.Sprintf(logFilePattern, appName)
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		funcLog("getLogFile").Fatalf("log file %s can't be created: %v", logFile, err)
	}
	return f
}
