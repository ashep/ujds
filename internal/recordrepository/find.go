package recordrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/ashep/ujds/internal/model"
	"github.com/ashep/ujds/internal/queryparser"
)

func (r *Repository) Find(
	ctx context.Context,
	index string,
	search string,
	since time.Time,
	cursor uint64,
	limit uint32,
) ([]model.Record, uint64, error) {
	if err := r.indexNameValidator.Validate(index); err != nil {
		return nil, 0, err //nolint:wrapcheck // ok
	}

	q := `SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at, r.touched_at FROM record r
		LEFT JOIN record_log l ON r.log_id = l.id
		LEFT JOIN index i ON r.index_id = i.id
		WHERE `
	qArgs := []any{}

	if search != "" {
		pq, err := queryparser.Parse(search)
		if err != nil {
			return nil, 0, fmt.Errorf("search query: %w", err)
		}

		qArgs = pq.Args()
		ql := len(qArgs)

		q += pq.String("r.data", 1)
		q += fmt.Sprintf(
			` AND i.name=$%d AND r.updated_at >= $%d AND l.id > $%d ORDER BY l.id LIMIT $%d`,
			ql+1, ql+2, ql+3, ql+4, //nolint:gomnd // ok
		)

		qArgs = append(qArgs, index, since, cursor, limit+1)
	} else {
		q += `i.name=$1 AND r.updated_at >= $2 AND l.id > $3 ORDER BY l.id LIMIT $4`

		qArgs = append(qArgs, index, since, cursor, limit+1)
	}

	rows, err := r.db.QueryContext(ctx, q, qArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("db query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	records := make([]model.Record, 0)
	recID, indexID, logID, data, crAt, upAt, tcAt := "", uint64(0), uint64(0), "", time.Time{}, time.Time{}, time.Time{}

	for rows.Next() {
		if err := rows.Scan(&recID, &indexID, &logID, &data, &crAt, &upAt, &tcAt); err != nil {
			return nil, 0, fmt.Errorf("db scan: %w", err)
		}

		records = append(records, model.Record{
			ID:        recID,
			IndexID:   indexID,
			Rev:       logID,
			Data:      data,
			CreatedAt: crAt,
			UpdatedAt: upAt,
			TouchedAt: tcAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("db rows iteration: %w", err)
	}

	newCursor := uint64(0)
	if len(records) > int(limit) {
		newCursor = records[limit-1].Rev
		records = records[:limit]
	}

	return records, newCursor, nil
}
