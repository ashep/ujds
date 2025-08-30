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

func TestIndexRepository_List(tt *testing.T) {
	tt.Run("DbQueryError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		dbm.
			ExpectQuery("SELECT .+ FROM index").
			WillReturnError(errors.New("theQueryError"))

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		_, err = repo.List(context.Background())

		assert.EqualError(t, err, "db query: theQueryError")
	})

	tt.Run("DbRowsIterationError", func(t *testing.T) {
		nameValidator := &stringValidatorMock{}

		db, dbm, err := sqlmock.New()
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "schema", "created_at", "updated_at"}).
			AddRow(123, "indexName", "{}", time.Unix(234, 0), time.Unix(345, 0)).
			RowError(0, errors.New("theRowError"))

		dbm.
			ExpectQuery("SELECT .+ FROM index").
			WillReturnRows(rows)

		repo := indexrepository.New(db, nameValidator, zerolog.Nop())
		_, err = repo.List(context.Background())

		assert.EqualError(t, err, "db rows iteration: theRowError")
	})
}
