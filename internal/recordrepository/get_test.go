package recordrepository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/recordrepository"
)

func TestRecordRepository_Get(tt *testing.T) {
	tt.Run("IndexNameValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return fmt.Errorf("theIndexNameValidationError")
		}

		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		_, err = repo.Get(context.Background(), "theIndexName", "theRecordID")
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

		jsonValidator := &jsonValidatorMock{}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		_, err = repo.Get(context.Background(), "theIndexName", "theRecordID")
		require.EqualError(t, err, "theRecordIDValidationError")
	})

	tt.Run("DbNoRows", func(t *testing.T) {
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

		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT r.index_id, r.log_id, l.data, r.created_at, r.updated_at, r.touched_at
FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id
WHERE i.name=$1 AND r.id=$2 ORDER BY l.created_at DESC LIMIT 1`).
			WithArgs("theIndexName", "theRecordID").
			WillReturnRows(sqlmock.NewRows([]string{}))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		_, err = repo.Get(context.Background(), "theIndexName", "theRecordID")
		require.ErrorIs(t, err, apperrors.NotFoundError{Subj: "record"})
	})

	tt.Run("DbRowScanError", func(t *testing.T) {
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

		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT r.index_id, r.log_id, l.data, r.created_at, r.updated_at, r.touched_at
FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id
WHERE i.name=$1 AND r.id=$2 ORDER BY l.created_at DESC LIMIT 1`).
			WithArgs("theIndexName", "theRecordID").
			WillReturnError(errors.New("theSQLError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		_, err = repo.Get(context.Background(), "theIndexName", "theRecordID")
		require.EqualError(t, err, "db scan: theSQLError")
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

		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"r.index_id", "r.log_id", "l.data", "r.created_at", "r.updated_at", "r.touched_at"})
		rows.AddRow(uint64(123), uint64(234), "theData", time.Unix(111, 0), time.Unix(112, 0), time.Unix(113, 0))
		dbm.
			ExpectQuery(`SELECT r.index_id, r.log_id, l.data, r.created_at, r.updated_at, r.touched_at
FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id
WHERE i.name=$1 AND r.id=$2 ORDER BY l.created_at DESC LIMIT 1`).
			WithArgs("theIndexName", "theRecordID").
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())

		rec, err := repo.Get(context.Background(), "theIndexName", "theRecordID")
		require.NoError(t, err)

		assert.Equal(t, "theRecordID", rec.ID)
		assert.Equal(t, uint64(123), rec.IndexID)
		assert.Equal(t, uint64(234), rec.Rev)
		assert.Equal(t, "theData", rec.Data)
		assert.Equal(t, time.Unix(111, 0), rec.CreatedAt)
		assert.Equal(t, time.Unix(112, 0), rec.UpdatedAt)
		assert.Equal(t, time.Unix(113, 0), rec.TouchedAt)
	})
}
