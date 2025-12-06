package recordrepo_test

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

	"github.com/ashep/ujds/internal/recordrepo"
)

func TestRecordRepository_History(tt *testing.T) {
	tt.Run("IndexNameValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndexName", s)
			return fmt.Errorf("theIndexNameValidationError")
		}

		recordIDValidator := &stringValidatorMock{}

		db, _, err := sqlmock.New()
		require.NoError(t, err)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())
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

		db, _, err := sqlmock.New()
		require.NoError(t, err)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())
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

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT`).
			WillReturnError(errors.New("theDbQueryError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())
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

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT`).
			WillReturnRows(rows)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())
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

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT`).
			WillReturnRows(sqlmock.NewRows([]string{}))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())
		res, cur, err := repo.History(context.Background(), "theIndexName", "theRecordID", time.Unix(0, 0), 0, 0)
		require.NoError(t, err)
		assert.Len(t, res, 0)
		assert.Equal(t, uint64(0), cur)
	})
}
