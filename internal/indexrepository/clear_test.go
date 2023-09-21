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

func TestRepository_Clear(tt *testing.T) {
	tt.Parallel()

	tt.Run("EmptyName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Clear(context.Background(), "")

		assert.ErrorIs(t, err, apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp ^[a-zA-Z0-9_/-]{1,255}$"})
	})

	tt.Run("InvalidName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Clear(context.Background(), "the n@me")

		assert.ErrorIs(t, err, apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp ^[a-zA-Z0-9_/-]{1,255}$"})
	})

	tt.Run("BeginTxError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin().WillReturnError(errors.New("theBeginTxError"))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		assert.EqualError(t, err, "failed to begin transaction: theBeginTxError")
	})

	tt.Run("ExecDeleteRecordsError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()

		dbm.ExpectExec(`DELETE FROM record WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnError(errors.New("theDeleteRecordsError"))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		assert.EqualError(t, err, "failed to delete records: theDeleteRecordsError")
	})

	tt.Run("ExecDeleteRecordLogsError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()

		dbm.ExpectExec(`DELETE FROM record WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbm.ExpectExec(`DELETE FROM record_log WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnError(errors.New("theDeleteRecordLogsError"))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		assert.EqualError(t, err, "failed to delete record logs: theDeleteRecordLogsError")
	})

	tt.Run("CommitError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()

		dbm.ExpectExec(`DELETE FROM record WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbm.ExpectExec(`DELETE FROM record_log WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbm.ExpectCommit().WillReturnError(errors.New("theCommitError"))

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		assert.EqualError(t, err, "db commit failed: theCommitError")
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()

		dbm.ExpectExec(`DELETE FROM record WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbm.ExpectExec(`DELETE FROM record_log WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbm.ExpectCommit()

		repo := indexrepository.New(db, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		require.NoError(t, err)
	})
}
