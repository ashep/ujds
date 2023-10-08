package recordrepository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/recordrepository"
)

func TestRecordRepository_History(tt *testing.T) {
	tt.Run("IndexNameValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return fmt.Errorf("theIndexNameValidationError")
		}

		recordIDValidator := &stringValidatorMock{}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		_, _, err = repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 0)
		require.EqualError(t, err, "theIndexNameValidationError")
	})

	tt.Run("RecordIDValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theRecordID", s)
			return fmt.Errorf("theRecordIDValidationError")
		}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		_, _, err = repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 0)
		require.EqualError(t, err, "theRecordIDValidationError")
	})

	tt.Run("DbQueryError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theRecordID", s)
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT id, index_id, data, created_at FROM record_log
WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1) AND record_id=$2 ORDER BY id DESC`).
			WithArgs("theIndexName", "theRecordID").
			WillReturnError(errors.New("theDbQueryError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		_, _, err = repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 0)
		require.EqualError(t, err, "db query: theDbQueryError")
	})

	tt.Run("DbRowsIterationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theRecordID", s)
			return nil
		}

		rows := sqlmock.NewRows([]string{"id", "index_id", "data", "created_at"}).
			RowError(0, errors.New("theRowError"))
		rows.AddRow(123, 234, `{}`, time.Time{})

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT id, index_id, data, created_at
FROM record_log WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1) AND record_id=$2 ORDER BY id DESC`).
			WithArgs("theIndexName", "theRecordID").
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		_, _, err = repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 0)
		require.EqualError(t, err, "db rows iteration: theRowError")
	})

	tt.Run("NoRows", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theRecordID", s)
			return nil
		}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT id, index_id, data, created_at
FROM record_log WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1) AND record_id=$2 ORDER BY id DESC`).
			WithArgs("theIndexName", "theRecordID").
			WillReturnRows(sqlmock.NewRows([]string{}))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		res, cur, err := repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 0)
		require.NoError(t, err)
		assert.Len(t, res, 0)
		assert.Equal(t, uint64(0), cur)
	})

	tt.Run("Ok", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theRecordID", s)
			return nil
		}

		rows := sqlmock.NewRows([]string{"id", "index_id", "data", "created_at"})
		rows.AddRow(123, 234, `{"foo":"bar"}`, time.Time{})

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT id, index_id, data, created_at
FROM record_log WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1) AND record_id=$2 ORDER BY id DESC`).
			WithArgs("theIndexName", "theRecordID").
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		res, cur, err := repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 0)
		require.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "theRecordID", res[0].ID)
		assert.Equal(t, uint64(234), res[0].IndexID)
		assert.Equal(t, uint64(123), res[0].Rev)
		assert.Equal(t, `{"foo":"bar"}`, res[0].Data)
		assert.Equal(t, uint64(0), cur)
	})

	tt.Run("OkLimit", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theRecordID", s)
			return nil
		}

		rows := sqlmock.NewRows([]string{"id", "index_id", "data", "created_at"})
		rows.AddRow(123, 234, `{"foo":"bar"}`, time.Time{})

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT id, index_id, data, created_at FROM record_log WHERE index_id=(SELECT id FROM index
WHERE name=$1 LIMIT 1) AND record_id=$2 ORDER BY id DESC LIMIT $3`).
			WithArgs("theIndexName", "theRecordID", int64(2)).
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		res, cur, err := repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 1)
		require.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "theRecordID", res[0].ID)
		assert.Equal(t, uint64(234), res[0].IndexID)
		assert.Equal(t, uint64(123), res[0].Rev)
		assert.Equal(t, `{"foo":"bar"}`, res[0].Data)

		assert.Equal(t, uint64(0), cur)
	})

	tt.Run("OkLimitOffset", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theRecordID", s)
			return nil
		}

		rows := sqlmock.NewRows([]string{"id", "index_id", "data", "created_at"})
		rows.AddRow(123, 234, `{"foo":"bar"}`, time.Time{})

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT id, index_id, data, created_at FROM record_log
WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1) AND record_id=$2 AND id<$3 ORDER BY id DESC LIMIT $4`).
			WithArgs("theIndexName", "theRecordID", uint64(123), int64(2)).
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		res, cur, err := repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 123, 1)
		require.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "theRecordID", res[0].ID)
		assert.Equal(t, uint64(234), res[0].IndexID)
		assert.Equal(t, uint64(123), res[0].Rev)
		assert.Equal(t, `{"foo":"bar"}`, res[0].Data)

		assert.Equal(t, uint64(0), cur)
	})

	tt.Run("OkMoreResults", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		recordIDValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theRecordID", s)
			return nil
		}

		rows := sqlmock.NewRows([]string{"id", "index_id", "data", "created_at"})
		rows.AddRow(123, 234, `{"foo":"bar"}`, time.Time{})
		rows.AddRow(345, 456, `{"foo2":"bar2"}`, time.Time{})

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT id, index_id, data, created_at FROM record_log WHERE index_id=(SELECT id FROM index
WHERE name=$1 LIMIT 1) AND record_id=$2 ORDER BY id DESC LIMIT $3`).
			WithArgs("theIndexName", "theRecordID", int64(2)).
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		res, cur, err := repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 1)
		require.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "theRecordID", res[0].ID)
		assert.Equal(t, uint64(234), res[0].IndexID)
		assert.Equal(t, uint64(123), res[0].Rev)
		assert.Equal(t, `{"foo":"bar"}`, res[0].Data)

		assert.Equal(t, uint64(123), cur)
	})
}
