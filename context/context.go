package context

import (
	"context"

	"github.com/sirupsen/logrus"
)

type indexContext int

const (
	loggerKey = iota
)

var (
	baseLogger logrus.FieldLogger
)

func init() {
	baseLogger = logrus.New().WithFields(nil)
}

func SetBaseLogger(logger logrus.FieldLogger) {
	baseLogger = logger
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}
	ctx = logWithField(ctx, "traceID", traceID)
	return ctx
}

func logWithField(ctx context.Context, key string, value interface{}) context.Context {
	log := Logger(ctx).WithField(key, value)
	return context.WithValue(ctx, loggerKey, log)
}

func Logger(ctx context.Context) logrus.FieldLogger {
	log, ok := ctx.Value(loggerKey).(*logrus.Entry)
	if !ok {
		return baseLogger
	}
	return log
}

func CheckDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
	return false
}
