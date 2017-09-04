package main

import (
	"context"
	"net/http"
	"time"

	pocontext "github.com/disiqueira/PoContext/context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func main() {
	router := chi.NewRouter()

	db, err := NewDB("")
	if err != nil {
		logrus.Fatal(err)
	}

	router.Use(
		pocontext.WithDatabase(db),
	)
	router.Route("/", func(r chi.Router) {
		router.Use(
			pocontext.WithTraceID(),
		)
		r.Get("/", NewGetHandler().ServeHTTP)
	})

	http.ListenAndServe(":8080", router)
}

type GetHandler struct {
}

func NewGetHandler() *GetHandler {
	return &GetHandler{}
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	someOtherFunc(r.Context())
}

func someOtherFunc(ctx context.Context) {
	log := pocontext.Logger(ctx)
	log.Info("Starting func")

	db, err := pocontext.DB(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	log.Infof("%+v", db)

	time.Sleep(time.Second * 5)

	log.Info("Finishing func")
}
