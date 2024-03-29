package indexrepository

import (
	"context"
	"fmt"

	"github.com/ashep/ujds/internal/model"
)

func (r *Repository) List(ctx context.Context) ([]model.Index, error) {
	q := "SELECT id, name, title, schema, created_at, updated_at FROM index"

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	res := make([]model.Index, 0)

	for rows.Next() {
		idx := model.Index{}
		if err := rows.Scan(&idx.ID, &idx.Name, &idx.Title, &idx.Schema, &idx.CreatedAt, &idx.UpdatedAt); err != nil {
			return nil, fmt.Errorf("db scan: %w", err)
		}

		res = append(res, idx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db rows iteration: %w", err)
	}

	return res, nil
}
