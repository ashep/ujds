package testapp

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq" // it's ok in tests
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/migration"
)

type Index struct {
	ID        int
	Name      string
	Schema    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TestDB struct {
	db *sql.DB
}

func newDB(t *testing.T) *TestDB {
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

func (d *TestDB) GetIndices(t *testing.T) []Index {
	t.Helper()

	rows, err := d.db.Query(`SELECT id, name, schema, created_at, updated_at FROM index`)
	require.NoError(t, err)

	res := make([]Index, 0)

	for rows.Next() {
		idx := Index{}
		require.NoError(t, rows.Scan(&idx.ID, &idx.Name, &idx.Schema, &idx.CreatedAt, &idx.UpdatedAt))
		res = append(res, idx)
	}

	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close()) //nolint:sqlclosecheck // this is testing code

	return res
}

func (d *TestDB) InsertIndex(t *testing.T, name, schema string) {
	t.Helper()

	_, err := d.db.Exec("INSERT INTO index (name, schema) VALUES ($1, $2)", name, schema)
	require.NoError(t, err)
}
