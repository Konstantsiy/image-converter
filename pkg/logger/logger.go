// Package logger provides a logging system.
package logger

import (
	"context"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type key int

const (
	loggerKey              key = 1
	DefaultTimestampFormat     = "2006-01-02 15:04:05"
)

var defaultLogger = &log.Logger{
	Out: os.Stdout,
	Formatter: &log.TextFormatter{
		DisableColors:   false,
		TimestampFormat: DefaultTimestampFormat,
		FullTimestamp:   true,
	},
	Level: log.InfoLevel,
}

// ContextWithLogger adds a defaultLogger to the given context.
func ContextWithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, defaultLogger)
}

// FromContext returns the logger from the given context.
func FromContext(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerKey).(*log.Logger); ok && logger != nil {
		return logger
	}
	return defaultLogger
}

// CompleteRequest logs http requests details like URI, method, status code and duration.
func CompleteRequest(ctx context.Context, r *http.Request, duration time.Duration, statusCode int) {
	FromContext(ctx).WithFields(log.Fields{
		"uri":      r.RequestURI,
		"method":   r.Method,
		"duration": duration,
		"status":   statusCode,
	}).Info("request completed")
}
