package indexrepo

import (
	"database/sql"

	"github.com/rs/zerolog"
)

//go:generate moq -out mock_test.go -pkg indexrepo_test -skip-ensure . stringValidator

type stringValidator interface {
	Validate(s string) error
}

type Repository struct {
	db            *sql.DB
	nameValidator stringValidator
	l             zerolog.Logger
}

func New(db *sql.DB, nameValidator stringValidator, l zerolog.Logger) *Repository {
	return &Repository{
		db:            db,
		nameValidator: nameValidator,
		l:             l,
	}
}
