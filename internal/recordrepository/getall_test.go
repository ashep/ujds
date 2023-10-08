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

func TestRecordRepository_GetAll(tt *testing.T) {
	tt.Run("IndexNameValidationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return fmt.Errorf("theIndexNameValidationError")
		}

		recordIDValidator := &stringValidatorMock{}

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		_, _, err = repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
		require.EqualError(t, err, "theIndexNameValidationError")
	})

	tt.Run("DbQueryError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at 
FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id 
WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4`).
			WithArgs("theIndex", time.Unix(123, 0), 234, 346).
			WillReturnError(errors.New("theDbError"))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		_, _, err = repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
		assert.EqualError(t, err, "db query: theDbError")
	})

	tt.Run("DbRowsIterationError", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"r.id", "r.index_id", "r.log_id", "l.data", "r.created_at", "r.updated_at"}).
			RowError(0, errors.New("theRowError"))
		rows.AddRow("", 0, 0, "", time.Time{}, time.Time{})

		dbm.ExpectQuery(`SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at 
FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id 
WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4`).
			WithArgs("theIndex", time.Unix(123, 0), 234, 346).
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		_, _, err = repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
		assert.EqualError(t, err, "db rows iteration: theRowError")
	})

	tt.Run("DbNoRows", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery(`SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at 
FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id 
WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4`).
			WithArgs("theIndex", time.Unix(123, 0), 234, 346).
			WillReturnRows(sqlmock.NewRows([]string{}))

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		res, cur, err := repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
		assert.NoError(t, err)
		assert.Empty(t, res)
		assert.Zero(t, cur)
	})

	tt.Run("OkNoMoreResults", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"r.id", "r.index_id", "r.log_id", "l.data", "r.created_at", "r.updated_at"})
		rows.AddRow("theID1", 12, 23, "theData1", time.Unix(111, 0), time.Unix(222, 0))
		rows.AddRow("theID2", 34, 45, "theData2", time.Unix(333, 0), time.Unix(444, 0))

		dbm.ExpectQuery(`SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at 
FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id 
WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4`).
			WithArgs("theIndex", time.Unix(123, 0), 234, 346).
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		res, cur, err := repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
		require.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, uint64(0), cur)

		assert.Equal(t, "theID1", res[0].ID)
		assert.Equal(t, uint64(12), res[0].IndexID)
		assert.Equal(t, uint64(23), res[0].Rev)
		assert.Equal(t, "theData1", res[0].Data)
		assert.Equal(t, time.Unix(111, 0), res[0].CreatedAt)
		assert.Equal(t, time.Unix(222, 0), res[0].UpdatedAt)

		assert.Equal(t, "theID2", res[1].ID)
		assert.Equal(t, uint64(34), res[1].IndexID)
		assert.Equal(t, uint64(45), res[1].Rev)
		assert.Equal(t, "theData2", res[1].Data)
		assert.Equal(t, time.Unix(333, 0), res[1].CreatedAt)
		assert.Equal(t, time.Unix(444, 0), res[1].UpdatedAt)
	})

	tt.Run("OkMoreResults", func(t *testing.T) {
		indexNameValidator := &stringValidatorMock{}
		indexNameValidator.ValidateFunc = func(s string) error {
			assert.Equal(t, "theIndex", s)
			return nil
		}

		recordIDValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"r.id", "r.index_id", "r.log_id", "l.data", "r.created_at", "r.updated_at"})
		rows.AddRow("theID1", 12, 23, "theData1", time.Unix(111, 0), time.Unix(222, 0))
		rows.AddRow("theID2", 34, 45, "theData2", time.Unix(333, 0), time.Unix(444, 0))
		rows.AddRow("theID2", 45, 56, "theData2", time.Unix(555, 0), time.Unix(666, 0))

		dbm.ExpectQuery(`SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at 
FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id 
WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4`).
			WithArgs("theIndex", time.Unix(123, 0), 234, 3).
			WillReturnRows(rows)

		repo := recordrepository.New(db, indexNameValidator, recordIDValidator, zerolog.Nop())

		res, cur, err := repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 2)
		require.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, uint64(45), cur)

		assert.Equal(t, "theID1", res[0].ID)
		assert.Equal(t, uint64(12), res[0].IndexID)
		assert.Equal(t, uint64(23), res[0].Rev)
		assert.Equal(t, "theData1", res[0].Data)
		assert.Equal(t, time.Unix(111, 0), res[0].CreatedAt)
		assert.Equal(t, time.Unix(222, 0), res[0].UpdatedAt)

		assert.Equal(t, "theID2", res[1].ID)
		assert.Equal(t, uint64(34), res[1].IndexID)
		assert.Equal(t, uint64(45), res[1].Rev)
		assert.Equal(t, "theData2", res[1].Data)
		assert.Equal(t, time.Unix(333, 0), res[1].CreatedAt)
		assert.Equal(t, time.Unix(444, 0), res[1].UpdatedAt)
	})
}
