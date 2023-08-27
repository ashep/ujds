package recordrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"

	"github.com/ashep/ujds/internal/model"
)

func (r *Repository) GetAll(
	ctx context.Context,
	index string,
	since time.Time,
	cursor uint64,
	limit uint32,
) ([]model.Record, uint64, error) {
	if index == "" {
		return nil, 0, apperrors.InvalidArgError{Subj: "index name", Reason: "must not be empty"}
	}

	if limit == 0 || limit > 500 {
		limit = 500
	}

	q := `SELECT r.id, r.log_id, l.data, r.created_at, r.updated_at FROM record r
		LEFT JOIN record_log l ON r.log_id = l.id
		LEFT JOIN index i ON r.index_id = i.id
		WHERE i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4`

	rows, err := r.db.QueryContext(ctx, q, index, since, cursor, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("db query failed: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	rcs := make([]model.Record, 0)
	recID, logID, data, crAt, upAt := "", uint64(0), "", time.Time{}, time.Time{}

	for rows.Next() {
		if err := rows.Scan(&recID, &logID, &data, &crAt, &upAt); err != nil {
			return nil, 0, fmt.Errorf("db scan failed: %w", err)
		}

		rcs = append(rcs, model.Record{
			ID:        recID,
			Index:     index,
			Rev:       logID,
			Data:      data,
			CreatedAt: crAt,
			UpdatedAt: upAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("db rows iteration failed: %w", err)
	}

	return rcs, logID, nil
}
