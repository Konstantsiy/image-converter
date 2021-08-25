// Package logger provides a logging system.
package logger

import (
	"context"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	loggerKey              = 1
	DefaultTimestampFormat = "2006-01-02 15:04:05"
)

var ctxLogger = &log.Logger{
	Out: os.Stdout,
	Formatter: &log.TextFormatter{
		DisableColors:   false,
		TimestampFormat: DefaultTimestampFormat,
		FullTimestamp:   true,
	},
	Level: log.InfoLevel,
}

// ContextWithLogger adds a ctxLogger to the given context.
func ContextWithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, ctxLogger)
}

// LoggerFromContext returns the logger from the given context.
func LoggerFromContext(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerKey).(*log.Logger); ok && logger != nil {
		return logger
	}
	return ctxLogger
}

// CompleteRequest logs http requests details like URI, method, status code and duration.
func CompleteRequest(ctx context.Context, r *http.Request, duration time.Duration, statusCode int) {
	LoggerFromContext(ctx).WithFields(log.Fields{
		"uri":      r.RequestURI,
		"method":   r.Method,
		"duration": duration,
		"status":   statusCode,
	}).Info("request completed")
}

// Info logs message at Info level.
func Info(ctx context.Context, msg string) {
	LoggerFromContext(ctx).Infoln(msg)
}

// Error logs errors at Error level.
func Error(ctx context.Context, err error) {
	LoggerFromContext(ctx).Errorln(err)
}
