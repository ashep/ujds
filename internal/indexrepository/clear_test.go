package indexrepository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/indexrepository"
)

func TestIndexRepository_Clear(tt *testing.T) {
	tt.Run("NameValidatorError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return errors.New("theValidatorError")
		}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Clear(context.Background(), "")

		assert.EqualError(t, err, "theValidatorError")
	})

	tt.Run("BeginTxError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin().WillReturnError(errors.New("theBeginTxError"))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		assert.EqualError(t, err, "begin transaction: theBeginTxError")
	})

	tt.Run("ExecDeleteRecordsError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()

		dbm.ExpectExec(`DELETE FROM record WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnError(errors.New("theDeleteRecordsError"))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		assert.EqualError(t, err, "delete records: theDeleteRecordsError")
	})

	tt.Run("ExecDeleteRecordLogsError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()

		dbm.ExpectExec(`DELETE FROM record WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbm.ExpectExec(`DELETE FROM record_log WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`).
			WithArgs("theIndex").
			WillReturnError(errors.New("theDeleteRecordLogsError"))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		assert.EqualError(t, err, "delete record log: theDeleteRecordLogsError")
	})

	tt.Run("CommitError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

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

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		assert.EqualError(t, err, "db commit: theCommitError")
	})

	tt.Run("Ok", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}
		nameValidator.ValidateFunc = func(s string) error {
			return nil
		}

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

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		err = repo.Clear(context.Background(), "theIndex")

		require.NoError(t, err)
	})
}
