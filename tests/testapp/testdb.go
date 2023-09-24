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
	Title     sql.NullString
	Schema    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RecordLog struct {
	ID        int
	IndexID   int
	RecordID  string
	Data      string
	CreatedAt time.Time
}

type Record struct {
	ID        string
	IndexID   int
	LogID     int
	Checksum  []byte
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

func (d *TestDB) GetIndex(t *testing.T, name string) Index {
	t.Helper()

	row := d.db.QueryRow(`SELECT id, name, title, schema, created_at, updated_at FROM index WHERE name=$1`, name)
	require.NoError(t, row.Err())

	idx := Index{}
	require.NoError(t, row.Scan(&idx.ID, &idx.Name, &idx.Title, &idx.Schema, &idx.CreatedAt, &idx.UpdatedAt))

	return idx
}

func (d *TestDB) GetIndices(t *testing.T) []Index {
	t.Helper()

	rows, err := d.db.Query(`SELECT id, name, title, schema, created_at, updated_at FROM index`)
	require.NoError(t, err)

	res := make([]Index, 0)

	for rows.Next() {
		idx := Index{}
		require.NoError(t, rows.Scan(&idx.ID, &idx.Name, &idx.Title, &idx.Schema, &idx.CreatedAt, &idx.UpdatedAt))
		res = append(res, idx)
	}

	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close()) //nolint:sqlclosecheck // this is testing code

	return res
}

func (d *TestDB) InsertIndex(t *testing.T, name, title, schema string) {
	t.Helper()

	sqlTitle := sql.NullString{
		String: title,
		Valid:  title != "",
	}

	_, err := d.db.Exec("INSERT INTO index (name, title, schema) VALUES ($1, $2, $3)", name, sqlTitle, schema)
	require.NoError(t, err)
}

func (d *TestDB) GetRecordLogs(t *testing.T, index string) []RecordLog {
	t.Helper()

	rows, err := d.db.Query(`SELECT id, index_id, record_id, data, created_at FROM record_log
WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`, index)
	require.NoError(t, err)

	res := make([]RecordLog, 0)

	for rows.Next() {
		rec := RecordLog{}
		require.NoError(t, rows.Scan(&rec.ID, &rec.IndexID, &rec.RecordID, &rec.Data, &rec.CreatedAt))
		res = append(res, rec)
	}

	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close()) //nolint:sqlclosecheck // this is testing code

	return res
}

func (d *TestDB) GetRecords(t *testing.T, index string) []Record {
	t.Helper()

	rows, err := d.db.Query(`SELECT id, index_id, log_id, checksum, created_at, updated_at FROM record
WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`, index)
	require.NoError(t, err)

	res := make([]Record, 0)

	for rows.Next() {
		rec := Record{}
		require.NoError(t, rows.Scan(&rec.ID, &rec.IndexID, &rec.LogID, &rec.Checksum, &rec.CreatedAt, &rec.UpdatedAt))
		res = append(res, rec)
	}

	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close()) //nolint:sqlclosecheck // this is testing code

	return res
}
