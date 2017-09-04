package context

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type indexContext int

const (
	traceIDKey = iota
	dbKey
)

var (
	fieldLogger logrus.FieldLogger
)

func init() {
	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime: "@timestamp",
		logrus.FieldKeyMsg:  "message",
	}}

	fieldLogger = l.WithField("key", "val")
}

func WithTraceID() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			traceID := r.Header.Get("ot-tracer-traceid")
			if traceID != "" {
				ctx = context.WithValue(ctx, traceIDKey, traceID)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func WithDatabase(db *sql.DB) func(next http.Handler) http.Handler {
	return middleware.WithValue(dbKey, db)
}

func Logger(ctx context.Context) logrus.FieldLogger {
	logger := fieldLogger

	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		logger = logger.WithField("traceID", traceID)
	}

	return logger
}

func TraceID(ctx context.Context) (string, error) {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return "", errors.New("No TraceID in the context")
	}

	return traceID, nil
}

func DB(ctx context.Context) (*sql.DB, error) {
	db, ok := ctx.Value(dbKey).(*sql.DB)
	if !ok {
		return nil, errors.New("No Database in the context")
	}

	return db, nil
}
