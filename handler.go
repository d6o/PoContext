package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	pocontext "github.com/disiqueira/PoContext/context"
	"github.com/palantir/stacktrace"
)

type GetHandler struct {
	db   *sql.DB
	time int
}

func NewGetHandler(db *sql.DB, time int) *GetHandler {
	return &GetHandler{
		db:   db,
		time: time,
	}
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Gets the logger using the context.
	log := pocontext.Logger(r.Context())
	log.Debug("Starting ServeHTTP function")
	defer log.Debug("Finishing ServeHTTP function")

	// Pass the context to the next level
	if err := someOtherFunc(r.Context(), h.db, h.time); err != nil {
		log.Warn(err)
		return
	}

	log.Info("Everything went well")
}

func someOtherFunc(ctx context.Context, db *sql.DB, timeAux int) error {
	log := pocontext.Logger(ctx)
	log.Debugf("Starting someOtherFunc function (timeAux: %d)", timeAux)
	defer log.Debug("Finishing someOtherFunc function")

	// Just sleep to simulate a database operation.
	query := "SELECT pg_sleep($1)"
	rows, err := db.QueryContext(ctx, query, timeAux)
	if err != nil {
		log.Warn(err)
		return stacktrace.Propagate(err, "error executing SQL Query")
	}
	defer rows.Close()

	// Simulate some other long operation.
	log.Debug("Starting sleep")
	time.Sleep(time.Duration(timeAux) * time.Second)

	//Verify if the context still valid.
	if pocontext.CheckTimeout(ctx) {
		return stacktrace.Propagate(err, "no more time to continue with this context")
	}

	log.Info("Executed everything in someOtherFunc")
	return nil
}
