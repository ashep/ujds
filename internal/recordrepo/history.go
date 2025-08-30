package recordrepo

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) History( //nolint:cyclop // TODO: calculated cyclomatic complexity for function History is 12, max is 10
	ctx context.Context,
	index string,
	id string,
	since time.Time,
	cursor uint64,
	limit uint32,
) ([]Record, uint64, error) {
	if err := r.indexNameValidator.Validate(index); err != nil {
		return nil, 0, err //nolint:wrapcheck // ok
	}

	if err := r.recordIDValidator.Validate(id); err != nil {
		return nil, 0, err //nolint:wrapcheck // ok
	}

	q := `SELECT id, index_id, data, created_at FROM record_log
WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1) AND record_id=$2`
	args := []interface{}{index, id}

	if since.Unix() != 0 {
		args = append(args, since)
		q += fmt.Sprintf(" AND created_at>=$%d", len(args))
	}

	if cursor != 0 {
		args = append(args, cursor)
		q += fmt.Sprintf(" AND id<$%d", len(args))
	}

	q += " ORDER BY id DESC"

	if limit != 0 {
		args = append(args, limit+1)
		q += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("db query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	records := make([]Record, 0)

	for rows.Next() {
		rec := Record{ID: id}
		if err := rows.Scan(&rec.Rev, &rec.IndexID, &rec.Data, &rec.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("db scan: %w", err)
		}

		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("db rows iteration: %w", err)
	}

	newCursor := uint64(0)
	if limit > 0 && len(records) > int(limit) {
		newCursor = records[limit-1].Rev
		records = records[:limit]
	}

	return records, newCursor, nil
}
