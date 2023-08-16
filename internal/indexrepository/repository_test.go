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

func TestRepository_Upsert(tt *testing.T) {
	tt.Parallel()

	tt.Run("EmptyName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "", "")

		require.ErrorIs(t, err, apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp ^[a-zA-Z0-9_-]{1,64}$"})
	})

	tt.Run("InvalidName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "the n@me", "")

		require.ErrorIs(t, err, apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp ^[a-zA-Z0-9_-]{1,64}$"})
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

	tt.Run("OkEmptySchema", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectExec(`INSERT INTO index (name, schema) VALUES ($1, $2) ON CONFLICT (name)
DO UPDATE SET schema=$2, updated_at=now()`).
			WithArgs("theIndex", "{}").
			WillReturnResult(sqlmock.NewResult(123, 234))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Upsert(context.Background(), "theIndex", "")

		require.NoError(t, err)
	})
}

func TestRepository_Get(tt *testing.T) {
	tt.Parallel()

	tt.Run("EmptyName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		_, err = repo.Get(context.Background(), "")

		require.ErrorIs(t, err, apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp ^[a-zA-Z0-9_-]{1,64}$"})
	})

	tt.Run("InvalidName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		_, err = repo.Get(context.Background(), "the n@me")

		require.ErrorIs(t, err, apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp ^[a-zA-Z0-9_-]{1,64}$"})
	})

	tt.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT id, schema, created_at, updated_at FROM index WHERE name=$1`).
			WithArgs("theIndex").
			WillReturnError(sql.ErrNoRows)

		repo := indexrepository.New(db, zerolog.Nop())
		_, err = repo.Get(context.Background(), "theIndex")

		require.ErrorIs(t, err, apperrors.NotFoundError{Subj: "index"})
	})

	tt.Run("DBScanError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT id, schema, created_at, updated_at FROM index WHERE name=$1`).
			WithArgs("theIndex").
			WillReturnError(errors.New("theDBExecError"))

		repo := indexrepository.New(db, zerolog.Nop())
		_, err = repo.Get(context.Background(), "theIndex")

		require.EqualError(t, err, "db scan failed: theDBExecError")
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "schema", "created_at", "updated_at"})
		rows.AddRow(123, `{"foo":"bar"}`, time.Unix(234, 0), time.Unix(345, 0))

		dbm.
			ExpectQuery(`SELECT id, schema, created_at, updated_at FROM index WHERE name=$1`).
			WithArgs("theIndex").
			WillReturnRows(rows)

		repo := indexrepository.New(db, zerolog.Nop())
		res, err := repo.Get(context.Background(), "theIndex")

		require.NoError(t, err)
		assert.Equal(t, 123, res.ID)
		assert.Equal(t, "theIndex", res.Name)
		assert.Equal(t, []byte(`{"foo":"bar"}`), res.Schema)
		assert.Equal(t, time.Unix(234, 0), res.CreatedAt)
		assert.Equal(t, time.Unix(345, 0), res.UpdatedAt)
	})
}
