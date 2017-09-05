package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"

	pocontext "github.com/disiqueira/PoContext/context"
	pomiddleware "github.com/disiqueira/PoContext/middleware"
)

func main() {
	db, err := NewDB(DATABASE_DSN)
	if err != nil {
		logrus.Fatal(err)
	}

	pocontext.SetLogger(logger())

	router := chi.NewRouter()
	router.Use(
		middleware.Timeout(5 * time.Second),
	)

	router.Route("/", func(r chi.Router) {
		router.Use(
			pomiddleware.TraceID(),
			pomiddleware.TimeoutRecover(),
		)
		r.Get("/fast", NewGetHandler(db, 2).ServeHTTP)
		r.Get("/medium", NewGetHandler(db, 3).ServeHTTP)
		r.Get("/slow", NewGetHandler(db, 8).ServeHTTP)
	})

	http.ListenAndServe(":8080", router)
}
