package recordrepository

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type Repository struct {
	db *sql.DB
	l  zerolog.Logger
}

func New(db *sql.DB, l zerolog.Logger) *Repository {
	return &Repository{db: db, l: l}
}
