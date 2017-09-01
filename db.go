package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewDB(DSN string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
