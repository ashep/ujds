package indexrepository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/indexrepository"
)

func TestRepository_List(tt *testing.T) {
	tt.Parallel()

	tt.Run("DbQueryError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		dbm.
			ExpectQuery("SELECT id, name, schema, created_at, updated_at FROM index").
			WillReturnError(errors.New("theQueryError"))

		repo := indexrepository.New(db, zerolog.Nop())
		_, err = repo.List(context.Background())

		assert.EqualError(t, err, "db query failed: theQueryError")
	})

	tt.Run("DbRowsIterationError", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "schema", "created_at", "updated_at"}).
			AddRow(123, "indexName", "{}", time.Unix(234, 0), time.Unix(345, 0)).
			RowError(0, errors.New("theRowError"))

		dbm.
			ExpectQuery("SELECT id, name, schema, created_at, updated_at FROM index").
			WillReturnRows(rows)

		repo := indexrepository.New(db, zerolog.Nop())
		_, err = repo.List(context.Background())

		assert.EqualError(t, err, "db rows iteration failed: theRowError")
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		db, dbm, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "schema", "created_at", "updated_at"}).
			AddRow(123, "indexName1", "indexSchema1", time.Unix(234, 0), time.Unix(345, 0)).
			AddRow(321, "indexName2", "indexSchema2", time.Unix(432, 0), time.Unix(543, 0))

		dbm.
			ExpectQuery("SELECT id, name, schema, created_at, updated_at FROM index").
			WillReturnRows(rows)

		repo := indexrepository.New(db, zerolog.Nop())
		res, err := repo.List(context.Background())

		require.NoError(t, err)
		assert.Len(t, res, 2)

		assert.Equal(t, uint64(123), res[0].ID)
		assert.Equal(t, "indexName1", res[0].Name)
		assert.Equal(t, []byte("indexSchema1"), res[0].Schema)
		assert.Equal(t, time.Unix(234, 0), res[0].CreatedAt)
		assert.Equal(t, time.Unix(345, 0), res[0].UpdatedAt)

		assert.Equal(t, uint64(321), res[1].ID)
		assert.Equal(t, "indexName2", res[1].Name)
		assert.Equal(t, []byte("indexSchema2"), res[1].Schema)
		assert.Equal(t, time.Unix(432, 0), res[1].CreatedAt)
		assert.Equal(t, time.Unix(543, 0), res[1].UpdatedAt)
	})
}
