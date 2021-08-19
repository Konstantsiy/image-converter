package logger

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	loggerKey              = 2
	DefaultTimestampFormat = "2006-01-02 15:04:05"
)

var ctxLogger = &log.Logger{
	Formatter: &log.TextFormatter{
		DisableColors:   false,
		TimestampFormat: DefaultTimestampFormat,
		FullTimestamp:   true,
	},
	Level: log.InfoLevel,
}

func GetCtxLogger() *log.Logger {
	return ctxLogger
}

func AddToContext(ctx context.Context, l *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func GetFormContext(ctx context.Context) *log.Logger {
	if l, ok := ctx.Value(loggerKey).(*log.Logger); ok && l != nil {
		return l
	}
	return ctxLogger
}

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
