package indexrepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
)

type Repository struct {
	db *sql.DB
	l  zerolog.Logger
}

func New(db *sql.DB, l zerolog.Logger) *Repository {
	return &Repository{db: db, l: l}
}

func (a *Repository) Upsert(ctx context.Context, name, schema string) error {
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

func (a *Repository) Get(ctx context.Context, name string) (Index, error) {
	var (
		id        int
		schema    []byte
		createdAt time.Time
		updatedAt time.Time
	)

	q := `SELECT id, schema, created_at, updated_at FROM index WHERE name=$1`

	row := a.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&id, &schema, &createdAt, &updatedAt); errors.Is(err, sql.ErrNoRows) {
		return Index{}, apperrors.NotFoundError{Subj: "schema"}
	} else if err != nil {
		return Index{}, fmt.Errorf("db scan failed: %w", err)
	}

	return Index{ID: id, Name: name, Schema: schema, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}
