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
	if err := r.nameValidator.Validate(name); err != nil {
		return model.Index{}, err //nolint:wrapcheck // ok
	}

	idx := model.Index{Name: name}
	q := `SELECT id, title, schema, created_at, updated_at FROM index WHERE name=$1`

	row := r.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&idx.ID, &idx.Title, &idx.Schema, &idx.CreatedAt, &idx.UpdatedAt); errors.Is(err, sql.ErrNoRows) {
		return model.Index{}, apperrors.NotFoundError{Subj: "index"}
	} else if err != nil {
		return model.Index{}, fmt.Errorf("db scan: %w", err)
	}

	return idx, nil
}
