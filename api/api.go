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

	mux         sync.RWMutex
	schemaCache map[string]Schema
}

func New(cfg Config, db *sql.DB, l zerolog.Logger) *API {
	return &API{
		cfg: cfg,
		db:  db,
		l:   l,

		mux:         sync.RWMutex{},
		schemaCache: make(map[string]Schema),
	}
}
