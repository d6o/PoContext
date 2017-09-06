package context

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type indexContext int

const (
	traceIDKey = iota
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
	if traceID != "" {
		ctx = context.WithValue(ctx, traceIDKey, traceID)
	}
	return ctx
}

func Logger(ctx context.Context) logrus.FieldLogger {
	logger := baseLogger

	if id, err := traceID(ctx); err == nil {
		logger = logger.WithField("traceID", id)
	}

	return logger
}

func traceID(ctx context.Context) (string, error) {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return "", errors.New("No TraceID in the context")
	}

	return traceID, nil
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
