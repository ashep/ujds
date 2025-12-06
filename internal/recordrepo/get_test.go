package recordrepo_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/recordrepo"
)

func TestRecordRepository_Get(tt *testing.T) {
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

		db, _, err := sqlmock.New()
		require.NoError(t, err)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())
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

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT`).
			WillReturnRows(sqlmock.NewRows([]string{}))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

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

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.
			ExpectQuery(`SELECT`).
			WillReturnError(errors.New("theSQLError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		_, err = repo.Get(context.Background(), "theIndexName", "theRecordID")
		require.EqualError(t, err, "db scan: theSQLError")
	})

}
