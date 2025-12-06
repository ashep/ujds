package recordrepo

import (
	"database/sql"

	"github.com/rs/zerolog"
)

//go:generate moq -out mock_test.go -pkg recordrepo_test -skip-ensure . stringValidator jsonValidator

type stringValidator interface {
	Validate(s string) error
}

type Repository struct {
	db                 *sql.DB
	indexNameValidator stringValidator
	recordIDValidator  stringValidator
	l                  zerolog.Logger
}

func New(
	db *sql.DB,
	indexNameValidator stringValidator,
	recordIDValidator stringValidator,
	l zerolog.Logger,
) *Repository {
	return &Repository{
		db:                 db,
		indexNameValidator: indexNameValidator,
		recordIDValidator:  recordIDValidator,
		l:                  l,
	}
}
