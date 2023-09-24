package indexrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ashep/go-apperrors"

	"github.com/ashep/ujds/internal/model"
)

func (r *Repository) Get(ctx context.Context, name string) (model.Index, error) {
	if !r.nameRe.MatchString(name) {
		return model.Index{}, apperrors.InvalidArgError{
			Subj: "name", Reason: "must match the regexp " + r.nameRe.String(),
		}
	}

	idx := model.Index{Name: name}
	q := `SELECT id, title, schema, created_at, updated_at FROM index WHERE name=$1`

	row := r.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&idx.ID, &idx.Title, &idx.Schema, &idx.CreatedAt, &idx.UpdatedAt); errors.Is(err, sql.ErrNoRows) {
		return model.Index{}, apperrors.NotFoundError{Subj: "index"}
	} else if err != nil {
		return model.Index{}, fmt.Errorf("db scan failed: %w", err)
	}

	return idx, nil
}
