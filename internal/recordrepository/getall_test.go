package recordrepository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/recordrepository"
)

func TestRepository_GetAll(tt *testing.T) {
	tt.Parallel()

	tt.Run("EmptyIndexName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		repo := recordrepository.New(db, zerolog.Nop())

		_, _, err = repo.GetAll(context.Background(), "", time.Time{}, 0, 0)
		assert.EqualError(t, err, "invalid index name: must not be empty")
	})

	tt.Run("DbQueryError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery("SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4").
			WithArgs("theIndex", time.Unix(123, 0), 234, 345).
			WillReturnError(errors.New("theDbError"))

		repo := recordrepository.New(db, zerolog.Nop())

		_, _, err = repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
		assert.EqualError(t, err, "db query failed: theDbError")
	})

	tt.Run("DbRowsError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"r.id", "r.index_id", "r.log_id", "l.data", "r.created_at", "r.updated_at"}).
			RowError(0, errors.New("theRowError"))
		rows.AddRow("", 0, 0, "", time.Time{}, time.Time{})

		dbm.ExpectQuery("SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4").
			WithArgs("theIndex", time.Unix(123, 0), 234, 345).
			WillReturnRows(rows)

		repo := recordrepository.New(db, zerolog.Nop())

		_, _, err = repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
		assert.EqualError(t, err, "db rows iteration failed: theRowError")
	})

	tt.Run("DbNoRows", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.ExpectQuery("SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4").
			WithArgs("theIndex", time.Unix(123, 0), 234, 345).
			WillReturnRows(sqlmock.NewRows([]string{}))

		repo := recordrepository.New(db, zerolog.Nop())

		res, cur, err := repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
		assert.NoError(t, err)
		assert.Empty(t, res)
		assert.Zero(t, cur)
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"r.id", "r.index_id", "r.log_id", "l.data", "r.created_at", "r.updated_at"})
		rows.AddRow("theID1", 12, 23, "theData1", time.Unix(111, 0), time.Unix(222, 0))
		rows.AddRow("theID2", 34, 45, "theData2", time.Unix(333, 0), time.Unix(444, 0))

		dbm.ExpectQuery("SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at FROM record r LEFT JOIN record_log l ON r.log_id = l.id LEFT JOIN index i ON r.index_id = i.id WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4").
			WithArgs("theIndex", time.Unix(123, 0), 234, 345).
			WillReturnRows(rows)

		repo := recordrepository.New(db, zerolog.Nop())

		res, cur, err := repo.GetAll(context.Background(), "theIndex", time.Unix(123, 0), 234, 345)
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
