package indexrepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
)

type Repository struct {
	db         *sql.DB
	nameRegexp *regexp.Regexp
	l          zerolog.Logger
}

func New(db *sql.DB, l zerolog.Logger) *Repository {
	return &Repository{
		db:         db,
		nameRegexp: regexp.MustCompile("^[a-zA-Z0-9_-]{1,64}$"),
		l:          l,
	}
}

func (a *Repository) Upsert(ctx context.Context, name, schema string) error {
	if !a.nameRegexp.MatchString(name) {
		return apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp " + a.nameRegexp.String()}
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

func (a *Repository) Get(ctx context.Context, name string) (model.Index, error) {
	var (
		id        int
		schema    []byte
		createdAt time.Time
		updatedAt time.Time
	)

	if !a.nameRegexp.MatchString(name) {
		return model.Index{}, apperrors.InvalidArgError{
			Subj: "name", Reason: "must match the regexp " + a.nameRegexp.String(),
		}
	}

	q := `SELECT id, schema, created_at, updated_at FROM index WHERE name=$1`

	row := a.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&id, &schema, &createdAt, &updatedAt); errors.Is(err, sql.ErrNoRows) {
		return model.Index{}, apperrors.NotFoundError{Subj: "index"}
	} else if err != nil {
		return model.Index{}, fmt.Errorf("db scan failed: %w", err)
	}

	return model.Index{ID: id, Name: name, Schema: schema, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}
