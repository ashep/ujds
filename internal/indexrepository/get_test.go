package indexrepository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/indexrepository"
)

func TestIndexRepository_Get(tt *testing.T) {
	tt.Run("NameValidatorError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return errors.New("theValidatorError")
		}

		db, _, err := sqlmock.New()
		require.NoError(t, err)

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		_, err = repo.Get(context.Background(), "")

		assert.EqualError(t, err, "theValidatorError")
	})

	tt.Run("NotFound", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT .+ FROM index`).
			WillReturnError(sql.ErrNoRows)

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		_, err = repo.Get(context.Background(), "theIndex")

		require.ErrorIs(t, err, apperrors.NotFoundError{Subj: "index"})
	})

	tt.Run("DBScanError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT .+ FROM index`).
			WillReturnError(errors.New("theDBExecError"))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		_, err = repo.Get(context.Background(), "theIndex")

		require.EqualError(t, err, "db scan: theDBExecError")
	})
}
