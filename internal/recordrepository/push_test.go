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
	tt.Run("ZeroIndexIDError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 0, nil, []model.RecordUpdate{})
		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index id",
			Reason: "must not be zero",
		})
	})

	tt.Run("EmptyRecordsError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.RecordUpdate{})
		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "records",
			Reason: "must not be empty",
		})
	})

	tt.Run("DbBeginError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin().WillReturnError(errors.New("theBeginError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.RecordUpdate{{}})
		require.EqualError(t, err, "db begin: theBeginError")
	})

	tt.Run("DbPrepareSelectError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1").
			WillReturnError(errors.New("thePrepareSelectError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.RecordUpdate{{}})
		require.EqualError(t, err, "db prepare: thePrepareSelectError")
	})

	tt.Run("DbPrepareInsertRecordLogError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WillReturnError(errors.New("thePrepareInsertRecordLogError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.RecordUpdate{{}})
		require.EqualError(t, err, "db prepare: thePrepareInsertRecordLogError")
	})

	tt.Run("DbPrepareInsertRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()").
			WillReturnError(errors.New("thePrepareInsertRecordError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.RecordUpdate{{}})
		require.EqualError(t, err, "db prepare: thePrepareInsertRecordError")
	})

	tt.Run("RecordIDValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

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
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.RecordUpdate{
			{ID: "theRecordID"},
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

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.RecordUpdate{
			{ID: "theRecordID", Data: ""},
		})
		require.EqualError(t, err, "invalid record data: must not be empty")
	})

	tt.Run("InvalidRecordDataJSON", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

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
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.RecordUpdate{
			{ID: "theRecordID", Data: "{]"},
		})
		require.EqualError(t, err, "invalid record data: invalid json")
	})

	tt.Run("RecordDataValidationFailed", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

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
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.RecordUpdate{
			{ID: "theRecordID", Data: "{}"},
		})

		require.EqualError(t, err, "invalid record data: (root): foo is required")
	})

	tt.Run("DbSelectRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

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
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}).
			WillReturnError(errors.New("theSelectError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.RecordUpdate{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "db scan: theSelectError")
	})

	tt.Run("DbInsertRecordLogError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

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
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnError(errors.New("theInsertRecordLogError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.RecordUpdate{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "db query: theInsertRecordLogError")
	})

	tt.Run("DbInsertRecordError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}, `{"foo":"bar"}`).
			WillReturnError(errors.New("theInsertRecordError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.RecordUpdate{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "db query: theInsertRecordError")
	})

	tt.Run("DbCommitError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}, `{"foo":"bar"}`).
			WillReturnResult(sqlmock.NewResult(345, 1))
		dbm.ExpectCommit().WillReturnError(errors.New("theCommitError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.RecordUpdate{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "db commit: theCommitError")
	})

	tt.Run("Ok", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}, `{"foo":"bar"}`).
			WillReturnResult(sqlmock.NewResult(345, 1))
		dbm.ExpectCommit()

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.RecordUpdate{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.NoError(t, err)
	})

	tt.Run("OkRecordAlreadyExists", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, s, "theRecordID")
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		selectLogRows := sqlmock.NewRows([]string{"log_id"})
		selectLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum, data) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{42, 74, 253, 163, 63, 3, 243, 26, 87, 206, 45, 219, 142, 20, 185, 244, 0, 171, 251, 145, 9, 55, 102, 88, 54, 182, 123, 225, 119, 28, 103, 187}).
			WillReturnRows(selectLogRows)
		dbm.ExpectCommit()

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.RecordUpdate{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.NoError(t, err)
	})
}
