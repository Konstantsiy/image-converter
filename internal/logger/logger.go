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
	loggerKey              = 2
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

// AddToContext adds a ctxLogger to the given context.
func AddToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, ctxLogger)
}

// GetFormContext returns the logger from the given context.
func GetFormContext(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerKey).(*log.Logger); ok && logger != nil {
		return logger
	}
	return ctxLogger
}

// CompleteRequest logs http requests details like URI, method, status code and duration.
func CompleteRequest(ctx context.Context, r *http.Request, duration time.Duration, statusCode int) {
	GetFormContext(ctx).WithFields(log.Fields{
		"uri":      r.RequestURI,
		"method":   r.Method,
		"duration": duration,
		"status":   statusCode,
	}).Info("request completed")
}

// Info logs message at Info level.
func Info(ctx context.Context, msg string) {
	GetFormContext(ctx).Infoln(msg)
}

// Error logs errors at Error level.
func Error(ctx context.Context, err error) {
	GetFormContext(ctx).Errorln(err)
}

// Fatal logs errors at Fatal level.
func Fatal(ctx context.Context, err error) {
	GetFormContext(ctx).Fatalln(err)
}
