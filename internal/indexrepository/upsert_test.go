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

func TestIndexRepository_Upsert(tt *testing.T) {
	tt.Run("NameValidatorError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return errors.New("theValidatorError")
		}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Upsert(context.Background(), "", "", "")

		assert.EqualError(t, err, "theValidatorError")
	})

	tt.Run("InvalidSchema", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "", "{]")

		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "schema",
			Reason: "invalid character ']' looking for beginning of object key string",
		})
	})

	tt.Run("DBExecError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) ON CONFLICT (name)
DO UPDATE SET title=$2, schema=$3, updated_at=now()`).
			WithArgs("theIndex", "theTitle", "{}").
			WillReturnError(errors.New("theDBExecError"))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "theTitle", "{}")

		require.EqualError(t, err, "db query failed: theDBExecError")
	})

	tt.Run("Ok", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) ON CONFLICT (name)
DO UPDATE SET title=$2, schema=$3, updated_at=now()`).
			WithArgs("theIndex", "theTitle", "{}").
			WillReturnResult(sqlmock.NewResult(123, 234))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "theTitle", "{}")

		require.NoError(t, err)
	})

	tt.Run("OkEmptySchema", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) ON CONFLICT (name)
DO UPDATE SET title=$2, schema=$3, updated_at=now()`).
			WithArgs("theIndex", sql.NullString{String: "theTitle", Valid: true}, "{}").
			WillReturnResult(sqlmock.NewResult(123, 234))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "theTitle", "")

		require.NoError(t, err)
	})

	tt.Run("OkEmptyTitle", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) ON CONFLICT (name)
DO UPDATE SET title=$2, schema=$3, updated_at=now()`).
			WithArgs("theIndex", nil, "{}").
			WillReturnResult(sqlmock.NewResult(123, 234))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "", "{}")

		require.NoError(t, err)
	})
}
