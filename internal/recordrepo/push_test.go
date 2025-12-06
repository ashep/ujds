package recordrepo_test

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

	"github.com/ashep/ujds/internal/recordrepo"
)

//nolint:maintidx // this is the test
func TestRecordRepository_Push(tt *testing.T) {
	tt.Run("EmptyUpdates", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, _, err := sqlmock.New()
		require.NoError(t, err)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{})
		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "updates",
			Reason: "must not be empty",
		})
	})

	tt.Run("DbBeginError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin().WillReturnError(errors.New("theBeginError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{{}})
		require.EqualError(t, err, "db begin: theBeginError")
	})

	tt.Run("DbPrepareSelectError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record").
			WillReturnError(errors.New("thePrepareSelectError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{{}})
		require.EqualError(t, err, "prepare statements: get record by log id: thePrepareSelectError")
	})

	tt.Run("DbPrepareInsertRecordLogError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log").
			WillReturnError(errors.New("thePrepareInsertRecordLogError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{{}})
		require.EqualError(t, err, "prepare statements: insert record log: thePrepareInsertRecordLogError")
	})

	tt.Run("DbPrepareInsertRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`).
			WillReturnError(errors.New("thePrepareInsertRecordError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{{}})
		require.EqualError(t, err, "prepare statements: insert record: thePrepareInsertRecordError")
	})

	tt.Run("DbPrepareTouchRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record `)
		dbm.ExpectPrepare(`UPDATE record`).
			WillReturnError(errors.New("thePrepareTouchRecordError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{{}})
		require.EqualError(t, err, "prepare statements: update record touch time: thePrepareTouchRecordError")
	})

	tt.Run("ZeroIndexID", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 0, ID: "theRecordID"},
		})
		require.EqualError(t, err, "invalid record 0: zero index id")
	})

	tt.Run("RecordIDValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return errors.New("theRecordIDValidationError")
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID"},
		})
		require.EqualError(t, err, "theRecordIDValidationError")
	})

	tt.Run("EmptyRecordData", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: ""},
		})
		require.EqualError(t, err, "invalid record data: must not be empty")
	})

	tt.Run("DbGetRecordByChecksumError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		dbm.ExpectQuery("SELECT log_id FROM record").
			WillReturnError(errors.New("theSelectError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "get record by checksum scan: theSelectError")
	})

	tt.Run("DbInsertLogError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		dbm.ExpectQuery("SELECT log_id FROM record").
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log").
			WillReturnError(errors.New("theInsertLogError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "upsert record: insert log db query: theInsertLogError")
	})

	tt.Run("DbInsertRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		dbm.ExpectQuery("SELECT log_id FROM record").
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log").
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record").
			WillReturnError(errors.New("theInsertRecordError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "upsert record: insert record db query: theInsertRecordError")
	})

	tt.Run("DbCommitError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		dbm.ExpectQuery("SELECT log_id FROM record").
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log").
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record").
			WillReturnResult(sqlmock.NewResult(345, 1))
		dbm.ExpectCommit().
			WillReturnError(errors.New("theCommitError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "commit: theCommitError")
	})

	tt.Run("Ok", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		dbm.ExpectQuery("SELECT log_id FROM record").
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log").
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record").
			WillReturnResult(sqlmock.NewResult(345, 1))
		dbm.ExpectCommit()

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.NoError(t, err)
	})

	tt.Run("DbUpdateTouchedAtExecError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		selectLogRows := sqlmock.NewRows([]string{"log_id"})
		selectLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		dbm.ExpectQuery(`SELECT log_id FROM record`).
			WillReturnRows(selectLogRows)
		dbm.ExpectExec(`UPDATE record SET touched_at`).
			WillReturnError(errors.New("theDbUpdateTouchedAtError"))
		dbm.ExpectCommit()

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "touch record: db exec: theDbUpdateTouchedAtError")
	})

	tt.Run("OkRecordAlreadyExists", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		selectLogRows := sqlmock.NewRows([]string{"log_id"})
		selectLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record")
		dbm.ExpectPrepare("INSERT INTO record_log")
		dbm.ExpectPrepare(`INSERT INTO record`)
		dbm.ExpectPrepare(`UPDATE record`)

		dbm.ExpectQuery(`SELECT log_id FROM record`).
			WillReturnRows(selectLogRows)
		dbm.ExpectExec(`UPDATE record SET touched_at=now()`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		dbm.ExpectCommit()

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []recordrepo.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.NoError(t, err)
	})
}
