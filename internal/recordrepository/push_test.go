package recordrepository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/model"
	"github.com/ashep/ujds/internal/recordrepository"
)

//nolint:maintidx // this is test
func TestRepository_Push(tt *testing.T) {
	tt.Parallel()

	tt.Run("ZeroIndexIDError", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 0, nil, []model.Record{})
		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index id",
			Reason: "must not be zero",
		})
	})

	tt.Run("EmptyRecordsError", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.Record{})
		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "records",
			Reason: "must not be empty",
		})
	})

	tt.Run("DbBeginError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin().WillReturnError(errors.New("theBeginError"))

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.Record{{}})
		require.EqualError(t, err, "db begin failed: theBeginError")
	})

	tt.Run("DbPrepareSelectError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1").
			WillReturnError(errors.New("thePrepareSelectError"))

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.Record{{}})
		require.EqualError(t, err, "db prepare failed: thePrepareSelectError")
	})

	tt.Run("DbInsertRecordLogError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WillReturnError(errors.New("thePrepareInsertRecordLogError"))

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.Record{{}})
		require.EqualError(t, err, "db prepare failed: thePrepareInsertRecordLogError")
	})

	tt.Run("DbInsertRecordError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()").
			WillReturnError(errors.New("thePrepareInsertRecordError"))

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.Record{{}})
		require.EqualError(t, err, "db prepare failed: thePrepareInsertRecordError")
	})

	tt.Run("EmptyRecordID", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.Record{
			{ID: ""},
		})
		require.EqualError(t, err, "invalid record (0) id: must not be empty")
	})

	tt.Run("EmptyRecordData", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.Record{
			{ID: "theRecordID", Data: ""},
		})
		require.EqualError(t, err, "invalid record (0) data: must not be empty")
	})

	tt.Run("InvalidRecordDataJSON", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, nil, []model.Record{
			{ID: "theRecordID", Data: "{]"},
		})
		require.EqualError(t, err, "invalid record data (0): invalid json")
	})

	tt.Run("RecordDataValidationFailed", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.Record{
			{ID: "theRecordID", Data: "{}"},
		})

		require.EqualError(t, err, "invalid record data (0): (root): foo is required")
	})

	tt.Run("DbSelectRecordError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnError(errors.New("theSelectError"))

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.Record{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "db scan failed: theSelectError")
	})

	tt.Run("DbInsertRecordLogError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnError(errors.New("theInsertRecordLogError"))

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.Record{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "db query failed: theInsertRecordLogError")
	})

	tt.Run("DbInsertRecordError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnError(errors.New("theInsertRecordError"))

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.Record{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "db query failed: theInsertRecordError")
	})

	tt.Run("DbCommitError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnResult(sqlmock.NewResult(345, 1))
		dbm.ExpectCommit().WillReturnError(errors.New("theCommitError"))

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.Record{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.EqualError(t, err, "db commit failed: theCommitError")
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		inertRecLogRows := sqlmock.NewRows([]string{"id"})
		inertRecLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnError(sql.ErrNoRows)
		dbm.ExpectQuery("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id").
			WithArgs(123, "theRecordID", `{"foo":"bar"}`).
			WillReturnRows(inertRecLogRows)
		dbm.ExpectExec("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()").
			WithArgs("theRecordID", 123, 234, []uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnResult(sqlmock.NewResult(345, 1))
		dbm.ExpectCommit()

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.Record{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.NoError(t, err)
	})

	tt.Run("OkRecordAlreadyExists", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		selectLogRows := sqlmock.NewRows([]string{"log_id"})
		selectLogRows.AddRow(234)

		dbm.ExpectBegin()
		dbm.ExpectPrepare("SELECT log_id FROM record WHERE checksum=$1")
		dbm.ExpectPrepare("INSERT INTO record_log (index_id, record_id, data) VALUES ($1, $2, $3) RETURNING id")
		dbm.ExpectPrepare("INSERT INTO record (id, index_id, log_id, checksum) VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()")
		dbm.ExpectQuery("SELECT log_id FROM record WHERE checksum=$1").
			WithArgs([]uint8{39, 77, 102, 66, 27, 133, 73, 126, 76, 125, 169, 231, 237, 162, 55, 166, 120, 101, 74, 40, 80, 24, 84, 188, 108, 45, 60, 193, 197, 122, 79, 242}).
			WillReturnRows(selectLogRows)
		dbm.ExpectCommit()

		repo := recordrepository.New(db, zerolog.Nop())

		err = repo.Push(context.Background(), 123, []byte(`{"required":["foo"]}`), []model.Record{
			{ID: "theRecordID", Data: `{"foo":"bar"}`},
		})

		require.NoError(t, err)
	})
}
