package indexrepository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/indexrepository"
)

func TestRepository_Upsert(tt *testing.T) {
	tt.Parallel()

	tt.Run("EmptyName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "", "", "")

		require.ErrorIs(t, err, apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp ^[a-zA-Z0-9_/-]{1,255}$"})
	})

	tt.Run("InvalidName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "the n@me", "", "")

		require.ErrorIs(t, err, apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp ^[a-zA-Z0-9_/-]{1,255}$"})
	})

	tt.Run("InvalidSchema", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "", "{]")

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
			ExpectExec(`INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) ON CONFLICT (name)
DO UPDATE SET title=$2, schema=$3, updated_at=now()`).
			WithArgs("theIndex", "theTitle", "{}").
			WillReturnError(errors.New("theDBExecError"))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "theTitle", "{}")

		require.EqualError(t, err, "db query failed: theDBExecError")
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) ON CONFLICT (name)
DO UPDATE SET title=$2, schema=$3, updated_at=now()`).
			WithArgs("theIndex", "theTitle", "{}").
			WillReturnResult(sqlmock.NewResult(123, 234))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "theTitle", "{}")

		require.NoError(t, err)
	})

	tt.Run("OkEmptySchema", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) ON CONFLICT (name)
DO UPDATE SET title=$2, schema=$3, updated_at=now()`).
			WithArgs("theIndex", sql.NullString{String: "theTitle", Valid: true}, "{}").
			WillReturnResult(sqlmock.NewResult(123, 234))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "theTitle", "")

		require.NoError(t, err)
	})

	tt.Run("OkEmptyTitle", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) ON CONFLICT (name)
DO UPDATE SET title=$2, schema=$3, updated_at=now()`).
			WithArgs("theIndex", nil, "{}").
			WillReturnResult(sqlmock.NewResult(123, 234))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "", "{}")

		require.NoError(t, err)
	})
}
