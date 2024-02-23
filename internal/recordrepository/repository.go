package recordrepository

import (
	"database/sql"

	"github.com/rs/zerolog"
)

//go:generate moq -out mock_test.go -pkg recordrepository_test -skip-ensure . stringValidator jsonValidator

type stringValidator interface {
	Validate(s string) error
}

type jsonValidator interface {
	Validate(schema, data []byte) error
}

type Repository struct {
	db                 *sql.DB
	indexNameValidator stringValidator
	recordIDValidator  stringValidator
	jsonValidator      jsonValidator
	l                  zerolog.Logger
}

func New(
	db *sql.DB,
	indexNameValidator stringValidator,
	recordIDValidator stringValidator,
	jsonValidator jsonValidator,
	l zerolog.Logger,
) *Repository {
	return &Repository{
		db:                 db,
		indexNameValidator: indexNameValidator,
		recordIDValidator:  recordIDValidator,
		jsonValidator:      jsonValidator,
		l:                  l,
	}
}
