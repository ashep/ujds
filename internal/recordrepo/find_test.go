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

func TestRecordRepository_Find(tt *testing.T) {
	tt.Run("IndexNameValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return fmt.Errorf("theIndexNameValidationError")
		}

		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, _, err := sqlmock.New()
		require.NoError(t, err)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())
		req := recordrepo.FindRequest{
			Index:           "theIndex",
			Query:           "",
			Since:           time.Unix(123, 0),
			Cursor:          234,
			Limit:           345,
			NotTouchedSince: nil,
		}
		_, _, err = repo.Find(context.Background(), req)
		require.EqualError(t, err, "theIndexNameValidationError")
	})

	tt.Run("InvalidSearchQuery", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		jsonValidator := &jsonValidatorMock{}
		recordIDValidator := &stringValidatorMock{}

		db, _, err := sqlmock.New()
		require.NoError(t, err)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())
		req := recordrepo.FindRequest{
			Index:           "theIndex",
			Query:           "foo bar baz",
			Since:           time.Unix(123, 0),
			Cursor:          234,
			Limit:           345,
			NotTouchedSince: nil,
		}

		_, _, err = repo.Find(context.Background(), req)
		require.EqualError(t, err, "search query: operator expected at position 4: foo ")
	})

	tt.Run("DbQueryError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT`).
			WillReturnError(errors.New("theDbError"))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())
		req := recordrepo.FindRequest{
			Index:           "theIndex",
			Query:           "",
			Since:           time.Unix(123, 0),
			Cursor:          234,
			Limit:           345,
			NotTouchedSince: nil,
		}

		_, _, err = repo.Find(context.Background(), req)
		assert.EqualError(t, err, "db query: theDbError")
	})

	tt.Run("DbRowsIterationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"r.id", "r.index_id", "r.log_id", "l.data", "r.created_at", "r.updated_at"}).
			RowError(0, errors.New("theRowError"))
		rows.AddRow("", 0, 0, "", time.Time{}, time.Time{})

		dbm.ExpectQuery(`SELECT`).
			WillReturnRows(rows)

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())
		req := recordrepo.FindRequest{
			Index:           "theIndex",
			Query:           "",
			Since:           time.Unix(123, 0),
			Cursor:          234,
			Limit:           345,
			NotTouchedSince: nil,
		}

		_, _, err = repo.Find(context.Background(), req)
		assert.EqualError(t, err, "db rows iteration: theRowError")
	})

	tt.Run("DbNoRows", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}
		jsonValidator := &jsonValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT`).
			WillReturnRows(sqlmock.NewRows([]string{}))

		repo := recordrepo.New(db, indexNameValidator, recordIDValidator, jsonValidator, zerolog.Nop())
		req := recordrepo.FindRequest{
			Index:           "theIndex",
			Query:           "",
			Since:           time.Unix(123, 0),
			Cursor:          234,
			Limit:           345,
			NotTouchedSince: nil,
		}

		res, cur, err := repo.Find(context.Background(), req)
		assert.NoError(t, err)
		assert.Empty(t, res)
		assert.Zero(t, cur)
	})
}
