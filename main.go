package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	router := chi.NewRouter()

	db, err := NewDB("")
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("%+v", db)

	ctx = WithDatabase(ctx, db)

	router.Get("/", NewGetHandler(ctx).ServeHTTP)
	http.ListenAndServe(":8888", router)
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

	db, err := DB(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	log.Infof("%+v", db)

	time.Sleep(time.Second * 5)

	log.Info("Finishing func")
}
