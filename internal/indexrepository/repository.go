package indexrepository

import (
	"database/sql"
	"regexp"

	"github.com/rs/zerolog"
)

type Repository struct {
	db     *sql.DB
	nameRe *regexp.Regexp
	l      zerolog.Logger
}

func New(db *sql.DB, l zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		nameRe: regexp.MustCompile("^[a-zA-Z0-9_/-]{1,255}$"),
		l:      l,
	}
}
