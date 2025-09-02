package recordrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ashep/ujds/internal/searchquery"
)

type FindRequest struct {
	Index           string
	Query           string
	Since           time.Time
	Cursor          uint64
	Limit           uint32
	NotTouchedSince *time.Time
}

func (r *Repository) Find(ctx context.Context, req FindRequest) ([]Record, uint64, error) {
	if err := r.indexNameValidator.Validate(req.Index); err != nil {
		return nil, 0, err //nolint:wrapcheck // ok
	}

	q := `SELECT r.id, r.index_id, r.log_id, l.data, r.created_at, r.updated_at, r.touched_at FROM record r
		LEFT JOIN record_log l ON r.log_id = l.id
		LEFT JOIN index i ON r.index_id = i.id
		WHERE `
	qArgs := []any{}

	if req.Query != "" {
		pq, err := searchquery.Parse(req.Query)
		if err != nil {
			return nil, 0, fmt.Errorf("search query: %w", err)
		}

		qArgs = pq.Args()
		ql := len(qArgs)

		q += pq.String("r.data", 1)
		q += fmt.Sprintf(` AND i.name=$%d AND r.updated_at >= $%d AND l.id > $%d`, ql+1, ql+2, ql+3)
		qArgs = append(qArgs, req.Index, req.Since, req.Cursor)
	} else {
		q += `i.name=$1 AND r.updated_at >= $2 AND l.id > $3`
		qArgs = append(qArgs, req.Index, req.Since, req.Cursor)
	}

	if req.NotTouchedSince != nil {
		q += fmt.Sprintf(` AND r.touched_at < $%d`, len(qArgs)+1)
		qArgs = append(qArgs, req.NotTouchedSince)
	}

	q += fmt.Sprintf(` ORDER BY l.id LIMIT $%d`, len(qArgs)+1)
	qArgs = append(qArgs, req.Limit+1)

	rows, err := r.db.QueryContext(ctx, q, qArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("db query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	records := make([]Record, 0)
	recID, indexID, logID, data, crAt, upAt, tcAt := "", uint64(0), uint64(0), "", time.Time{}, time.Time{}, time.Time{}

	for rows.Next() {
		if err := rows.Scan(&recID, &indexID, &logID, &data, &crAt, &upAt, &tcAt); err != nil {
			return nil, 0, fmt.Errorf("db scan: %w", err)
		}

		records = append(records, Record{
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
	if len(records) > int(req.Limit) {
		newCursor = records[req.Limit-1].Rev
		records = records[:req.Limit]
	}

	return records, newCursor, nil
}
