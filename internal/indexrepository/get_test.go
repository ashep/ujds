package indexrepository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

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

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
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

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT id, title, schema, created_at, updated_at FROM index WHERE name=$1`).
			WithArgs("theIndex").
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

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT id, title, schema, created_at, updated_at FROM index WHERE name=$1`).
			WithArgs("theIndex").
			WillReturnError(errors.New("theDBExecError"))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		_, err = repo.Get(context.Background(), "theIndex")

		require.EqualError(t, err, "db scan: theDBExecError")
	})

	tt.Run("Ok", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "title", "schema", "created_at", "updated_at"})
		rows.AddRow(123, "theTitle", `{"foo":"bar"}`, time.Unix(234, 0), time.Unix(345, 0))

		dbm.
			ExpectQuery(`SELECT id, title, schema, created_at, updated_at FROM index WHERE name=$1`).
			WithArgs("theIndex").
			WillReturnRows(rows)

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		res, err := repo.Get(context.Background(), "theIndex")

		require.NoError(t, err)
		assert.Equal(t, uint64(123), res.ID)
		assert.Equal(t, "theIndex", res.Name)
		assert.Equal(t, "theTitle", res.Title.String)
		assert.Equal(t, []byte(`{"foo":"bar"}`), res.Schema)
		assert.Equal(t, time.Unix(234, 0), res.CreatedAt)
		assert.Equal(t, time.Unix(345, 0), res.UpdatedAt)
	})
}
