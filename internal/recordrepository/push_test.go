package recordrepository_test

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

	"github.com/ashep/ujds/internal/model"
	"github.com/ashep/ujds/internal/recordrepository"
)

//nolint:maintidx // this is the test
func TestRecordRepository_Push(tt *testing.T) {
	tt.Run("EmptyUpdates", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{})
		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "updates",
			Reason: "must not be empty",
		})
	})

	tt.Run("DbBeginError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin().WillReturnError(errors.New("theBeginError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{{}})
		require.EqualError(t, err, "db begin: theBeginError")
	})

	tt.Run("DbPrepareSelectError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1").
			WillReturnError(errors.New("thePrepareSelectError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{{}})
		require.EqualError(t, err, "prepare statements: get record by log id: thePrepareSelectError")
	})

	tt.Run("DbPrepareInsertRecordLogError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WillReturnError(errors.New("thePrepareInsertRecordLogError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{{}})
		require.EqualError(t, err, "prepare statements: insert record log: thePrepareInsertRecordLogError")
	})

	tt.Run("DbPrepareInsertRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`).
			WillReturnError(errors.New("thePrepareInsertRecordError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{{}})
		require.EqualError(t, err, "prepare statements: insert record: thePrepareInsertRecordError")
	})

	tt.Run("DbPrepareTouchRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`).
			WillReturnError(errors.New("thePrepareTouchRecordError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{{}})
		require.EqualError(t, err, "prepare statements: update record touch time: thePrepareTouchRecordError")
	})

	tt.Run("ZeroIndexID", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
			{IndexID: 0, ID: "theRecordID"},
		})
		require.EqualError(t, err, "invalid record 0: zero index id")
	})

	tt.Run("RecordIDValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return errors.New("theRecordIDValidationError")
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
			{IndexID: 123, ID: "theRecordID"},
		})
		require.EqualError(t, err, "theRecordIDValidationError")
	})

	tt.Run("EmptyRecordData", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: ""},
		})
		require.EqualError(t, err, "invalid record data: must not be empty")
	})

	tt.Run("JSONValidatorError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		jsonValidator := &jsonValidatorMock{}
		jsonValidator.ValidateFunc = func(schema []byte, data []byte) error {
			return errors.New("theJSONValidatorError")
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: "{]"},
		})
		require.EqualError(t, err, "invalid record data: theJSONValidatorError")
	})

	tt.Run("DbSelectRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		jsonValidator := &jsonValidatorMock{}
		jsonValidator.ValidateFunc = func(schema []byte, data []byte) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}).
			WillReturnError(errors.New("theSelectError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
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

		jsonValidator := &jsonValidatorMock{}
		jsonValidator.ValidateFunc = func(schema []byte, data []byte) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnError(errors.New("theInsertLogError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
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

		jsonValidator := &jsonValidatorMock{}
		jsonValidator.ValidateFunc = func(schema []byte, data []byte) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}, `{"foo":"bar"}`).
			WillReturnError(errors.New("theInsertRecordError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
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

		jsonValidator := &jsonValidatorMock{}
		jsonValidator.ValidateFunc = func(schema []byte, data []byte) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}, `{"foo":"bar"}`).
			WillReturnResult(sqlmock.NewResult(345, 1))
		dbm.ExpectCommit().WillReturnError(errors.New("theCommitError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
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

		jsonValidator := &jsonValidatorMock{}
		jsonValidator.ValidateFunc = func(schema []byte, data []byte) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}, `{"foo":"bar"}`).
			WillReturnResult(sqlmock.NewResult(345, 1))
		dbm.ExpectCommit()

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
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

		jsonValidator := &jsonValidatorMock{}
		jsonValidator.ValidateFunc = func(schema []byte, data []byte) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		selectLogRows := sqlmock.NewRows([]string{"log_id"})
		selectLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		dbm.ExpectQuery(`SELECT log_id FROM record WHERE checksum=$1`).
			WithArgs([]uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}).
			WillReturnRows(selectLogRows)
		dbm.ExpectExec(`UPDATE record SET touched_at=now() WHERE log_id=$1`).
			WithArgs(234).
			WillReturnError(errors.New("theDbUpdateTouchedAtError"))
		dbm.ExpectCommit()

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
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

		jsonValidator := &jsonValidatorMock{}
		jsonValidator.ValidateFunc = func(schema []byte, data []byte) error {
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		selectLogRows := sqlmock.NewRows([]string{"log_id"})
		selectLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare(`INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
		dbm.ExpectPrepare(`UPDATE record SET touched_at=now() WHERE log_id=$1`)

		dbm.ExpectQuery(`SELECT log_id FROM record WHERE checksum=$1`).
			WithArgs([]uint8{207, 14, 59, 238, 143, 117, 105, 162, 113, 60, 2, 24, 160, 174, 111, 40, 180, 35, 202, 226, 143, 106, 209, 59, 233, 175, 54, 219, 8, 181, 47, 149}).
			WillReturnRows(selectLogRows)
		dbm.ExpectExec(`UPDATE record SET touched_at=now() WHERE log_id=$1`).
			WithArgs(234).
			WillReturnResult(sqlmock.NewResult(1, 1))
		dbm.ExpectCommit()

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		err = repo.Push(context.Background(), []model.RecordUpdate{
			{IndexID: 123, ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.NoError(t, err)
	})
}
