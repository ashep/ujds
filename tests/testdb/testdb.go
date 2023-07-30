package testdb

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // it's ok in tests
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/migration"
)

type TestDB struct {
	db *sql.DB
}

func New(t *testing.T) *TestDB {
	t.Helper()

	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable")

	require.NoError(t, err)
	require.NoError(t, db.Ping())

	return &TestDB{db: db}
}

func (d *TestDB) Reset(t *testing.T) {
	t.Helper()

	require.NoError(t, migration.Down(d.db))
	require.NoError(t, migration.Up(d.db))
}
