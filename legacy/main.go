package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type indexContext int

const (
	traceIDKey = iota
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

func Logger(ctx context.Context) logrus.FieldLogger {
	logger := fieldLogger

	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		logger = logger.WithField("traceID", traceID)
	}

	return logger
}

type GetHandler struct {
	ctx context.Context
}

func NewGetHandler(ctx context.Context) *GetHandler {
	return &GetHandler{
		ctx: ctx,
	}
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := WithTraceID(h.ctx, uuid.NewV4().String())

	someOtherFunc(ctx)
}

func someOtherFunc(ctx context.Context) {
	log := Logger(ctx)
	log.Info("Starting func")

	time.Sleep(time.Second * 5)

	log.Info("Finishing func")
}

func main() {
	ctx := context.Background()
	router := chi.NewRouter()

	router.Get("/", NewGetHandler(ctx).ServeHTTP)
	http.ListenAndServe(":8888", router)
}
