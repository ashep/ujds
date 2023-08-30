package indexrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"

	"github.com/ashep/ujds/internal/model"
)

func (r *Repository) Get(ctx context.Context, name string) (model.Index, error) {
	var (
		id        uint64
		schema    []byte
		createdAt time.Time
		updatedAt time.Time
	)

	if !r.nameRe.MatchString(name) {
		return model.Index{}, apperrors.InvalidArgError{
			Subj: "name", Reason: "must match the regexp " + r.nameRe.String(),
		}
	}

	q := `SELECT id, schema, created_at, updated_at FROM index WHERE name=$1`

	row := r.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&id, &schema, &createdAt, &updatedAt); errors.Is(err, sql.ErrNoRows) {
		return model.Index{}, apperrors.NotFoundError{Subj: "index"}
	} else if err != nil {
		return model.Index{}, fmt.Errorf("db scan failed: %w", err)
	}

	return model.Index{ID: id, Name: name, Schema: schema, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}
