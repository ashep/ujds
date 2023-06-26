package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/ashep/ujds/internal/errs"
)

func (a *API) UpsertIndex(ctx context.Context, name, schema string) error {
	if name == "" {
		return errs.ErrEmptyArg{Subj: "name"}
	}
	if schema == "" {
		schema = "{}"
	}

	if err := json.Unmarshal([]byte(schema), &struct{}{}); err != nil {
		return errs.ErrInvalidArg{Subj: "schema", E: err}
	}

	q := `INSERT INTO index (name, schema) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET schema=$2, updated_at=now()`
	if _, err := a.db.ExecContext(ctx, q, name, schema); err != nil {
		return err
	}

	return nil
}

func (a *API) GetIndex(ctx context.Context, name string) (*Schema, error) {
	var (
		id        int
		data      []byte
		createdAt time.Time
		updatedAt time.Time
	)

	q := `SELECT id, schema, created_at, updated_at FROM index WHERE name=$1`
	row := a.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&id, &data, &createdAt, &updatedAt); errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound{Subj: "schema"}
	} else if err != nil {
		return nil, err
	}

	return &Schema{Id: id, Name: name, Data: data, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}
