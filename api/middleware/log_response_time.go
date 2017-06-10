package middleware

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

type logResponseTime struct {
	logger  *log.Logger
	handler http.Handler
}

func NewLogResponseTime(l *log.Logger, h http.Handler) http.Handler {
	return &logResponseTime{
		logger:  l,
		handler: h,
	}
}

func (logResponseTime *logResponseTime) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	started := time.Now()

	logResponseTime.handler.ServeHTTP(w, r)

	logResponseTime.logger.WithFields(log.Fields{
		"method":        r.Method,
		"path":          r.URL.Path,
		"response_time": time.Since(started),
	}).Info("Request has been processed")
}
