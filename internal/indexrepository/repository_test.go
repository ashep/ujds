package indexrepository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/indexrepository"
)

func TestAPI_UpsertIndex(tt *testing.T) {
	tt.Parallel()

	tt.Run("EmptyName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "", "")

		require.ErrorIs(t, err, apperrors.EmptyArgError{Subj: "name"})
		assert.EqualError(t, err, "name is empty")
	})

	tt.Run("InvalidSchema", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "{]")

		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "schema",
			Reason: "invalid character ']' looking for beginning of object key string",
		})
	})

	tt.Run("DBExecError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, schema) VALUES ($1, $2) ON CONFLICT (name)
DO UPDATE SET schema=$2, updated_at=now()`).
			WithArgs("theIndex", "{}").
			WillReturnError(errors.New("theDBExecError"))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "{}")

		require.EqualError(t, err, "db query failed: theDBExecError")
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, schema) VALUES ($1, $2) ON CONFLICT (name)
DO UPDATE SET schema=$2, updated_at=now()`).
			WithArgs("theIndex", "{}").
			WillReturnResult(sqlmock.NewResult(123, 234))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "{}")

		require.NoError(t, err)
	})
}
