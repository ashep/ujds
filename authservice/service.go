package authservice

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type Service struct {
	db  *sql.DB
	cfg Config
	l   zerolog.Logger
}

func New(db *sql.DB, cfg Config, l zerolog.Logger) *Service {
	return &Service{
		db:  db,
		cfg: cfg,
		l:   l,
	}
}
