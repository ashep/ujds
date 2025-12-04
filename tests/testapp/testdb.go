package testapp

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ashep/go-app/testpostgres"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq" // ok in tests
	"github.com/stretchr/testify/require"
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
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
	TouchedAt time.Time
}

type TestDB struct {
	DSN string

	t *testing.T
	d *sql.DB
}

func newDB(t *testing.T) *TestDB {
	t.Helper()

	tp := testpostgres.New(t)

	return &TestDB{
		DSN: tp.DSN(),

		t: t,
		d: stdlib.OpenDBFromPool(tp.DB()),
	}
}

func (d *TestDB) GetIndex(name string) Index {
	row := d.d.QueryRow(`SELECT id, name, title, schema, created_at, updated_at FROM index WHERE name=$1`, name)
	require.NoError(d.t, row.Err())

	idx := Index{}
	require.NoError(d.t, row.Scan(&idx.ID, &idx.Name, &idx.Title, &idx.Schema, &idx.CreatedAt, &idx.UpdatedAt))

	return idx
}

func (d *TestDB) GetIndices() []Index {
	rows, err := d.d.Query(`SELECT id, name, title, schema, created_at, updated_at FROM index`)
	require.NoError(d.t, err)

	res := make([]Index, 0)

	for rows.Next() {
		idx := Index{}
		require.NoError(d.t, rows.Scan(&idx.ID, &idx.Name, &idx.Title, &idx.Schema, &idx.CreatedAt, &idx.UpdatedAt))
		res = append(res, idx)
	}

	require.NoError(d.t, rows.Err())
	require.NoError(d.t, rows.Close()) //nolint:sqlclosecheck // this is testing code

	return res
}

func (d *TestDB) InsertIndex(name, title, schema string) {
	sqlTitle := sql.NullString{
		String: title,
		Valid:  title != "",
	}

	_, err := d.d.Exec("INSERT INTO index (name, title, schema) VALUES ($1, $2, $3)", name, sqlTitle, schema)
	require.NoError(d.t, err)
}

func (d *TestDB) GetRecordLogs(index string) []RecordLog {
	rows, err := d.d.Query(`SELECT id, index_id, record_id, data, created_at FROM record_log
WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`, index)
	require.NoError(d.t, err)

	res := make([]RecordLog, 0)

	for rows.Next() {
		rec := RecordLog{}
		require.NoError(d.t, rows.Scan(&rec.ID, &rec.IndexID, &rec.RecordID, &rec.Data, &rec.CreatedAt))
		res = append(res, rec)
	}

	require.NoError(d.t, rows.Err())
	require.NoError(d.t, rows.Close()) //nolint:sqlclosecheck // this is testing code

	return res
}

func (d *TestDB) GetRecords(index string) []Record {
	rows, err := d.d.Query(`SELECT id, index_id, log_id, checksum, data, created_at, updated_at, touched_at FROM record
WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`, index)
	require.NoError(d.t, err)

	res := make([]Record, 0)

	for rows.Next() {
		rec := Record{}
		require.NoError(d.t, rows.Scan(&rec.ID, &rec.IndexID, &rec.LogID, &rec.Checksum, &rec.Data, &rec.CreatedAt, &rec.UpdatedAt, &rec.TouchedAt))
		res = append(res, rec)
	}

	require.NoError(d.t, rows.Err())
	require.NoError(d.t, rows.Close()) //nolint:sqlclosecheck // this is testing code

	return res
}
