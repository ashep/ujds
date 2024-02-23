package recordrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ashep/go-apperrors"

	"github.com/ashep/ujds/internal/model"
)

// Get returns last version of a record.
func (r *Repository) Get(ctx context.Context, index, id string) (model.Record, error) {
	if err := r.indexNameValidator.Validate(index); err != nil {
		return model.Record{}, err //nolint:wrapcheck // ok
	}

	if err := r.recordIDValidator.Validate(id); err != nil {
		return model.Record{}, err //nolint:wrapcheck // ok
	}

	q := `SELECT r.index_id, r.log_id, l.data, r.created_at, r.updated_at, r.touched_at FROM record r
		LEFT JOIN record_log l ON r.log_id = l.id
		LEFT JOIN index i ON r.index_id = i.id
		WHERE i.name=$1 AND r.id=$2 ORDER BY l.created_at DESC LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, index, id)

	rec := model.Record{
		ID: id,
	}

	err := row.Scan(&rec.IndexID, &rec.Rev, &rec.Data, &rec.CreatedAt, &rec.UpdatedAt, &rec.TouchedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Record{}, apperrors.NotFoundError{Subj: "record"}
	} else if err != nil {
		return model.Record{}, fmt.Errorf("db scan: %w", err)
	}

	return rec, nil
}
