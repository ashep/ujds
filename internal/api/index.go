package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/xeipuuv/gojsonschema"

	"github.com/ashep/ujds/internal/apperrors"
)

type Index struct {
	ID        int
	Name      string
	Schema    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Index) Validate(data []byte) error {
	if len(s.Schema) == 0 || bytes.Equal(s.Schema, []byte("{}")) {
		return nil
	}

	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(s.Schema), gojsonschema.NewBytesLoader(data))
	if err != nil {
		return fmt.Errorf("schema validate failed: %w", err)
	}

	if !res.Valid() {
		return errors.New(res.Errors()[0].String())
	}

	return nil
}

func (a *API) UpsertIndex(ctx context.Context, name, schema string) error {
	if name == "" {
		return apperrors.EmptyArgError{Subj: "name"}
	}

	if schema == "" {
		schema = "{}"
	}

	if err := json.Unmarshal([]byte(schema), &struct{}{}); err != nil {
		return apperrors.InvalidArgError{Subj: "schema", Reason: err.Error()}
	}

	q := `INSERT INTO index (name, schema) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET schema=$2, updated_at=now()`
	if _, err := a.db.ExecContext(ctx, q, name, schema); err != nil {
		return fmt.Errorf("db query failed: %w", err)
	}

	return nil
}

func (a *API) GetIndex(ctx context.Context, name string) (*Index, error) {
	var (
		id        int
		schema    []byte
		createdAt time.Time
		updatedAt time.Time
	)

	q := `SELECT id, schema, created_at, updated_at FROM index WHERE name=$1`

	row := a.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&id, &schema, &createdAt, &updatedAt); errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NotFoundError{Subj: "schema"}
	} else if err != nil {
		return nil, fmt.Errorf("db scan failed: %w", err)
	}

	return &Index{ID: id, Name: name, Schema: schema, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}
