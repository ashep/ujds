package api

import (
	"database/sql"
	"sync"

	"github.com/rs/zerolog"
)

type API struct {
	cfg Config
	db  *sql.DB
	l   zerolog.Logger

	mux sync.RWMutex
}

func New(cfg Config, db *sql.DB, l zerolog.Logger) *API {
	return &API{
		cfg: cfg,
		db:  db,
		l:   l,
	}
}
