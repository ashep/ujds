package indexrepo

import (
	"context"
	"fmt"
)

func (r *Repository) List(ctx context.Context) ([]Index, error) {
	q := "SELECT id, name, title, created_at, updated_at FROM index"

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	res := make([]Index, 0)

	for rows.Next() {
		idx := Index{}
		if err := rows.Scan(&idx.ID, &idx.Name, &idx.Title, &idx.CreatedAt, &idx.UpdatedAt); err != nil {
			return nil, fmt.Errorf("db scan: %w", err)
		}

		res = append(res, idx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db rows iteration: %w", err)
	}

	return res, nil
}
