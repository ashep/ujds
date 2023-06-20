package api

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type API struct {
	db *sql.DB
	l  zerolog.Logger
}

func New(db *sql.DB, l zerolog.Logger) *API {
	return &API{
		db: db,
		l:  l,
	}
}
