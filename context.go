package main

import (
	"context"
	"database/sql"
	"errors"
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

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func WithDatabase(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
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
