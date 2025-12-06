package indexrepo

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *Repository) Upsert(ctx context.Context, name, title string) error {
	if err := r.nameValidator.Validate(name); err != nil {
		return err //nolint:wrapcheck // ok
	}

	sqlTitle := sql.NullString{
		String: title,
		Valid:  title != "",
	}

	q := `INSERT INTO index (name, title) VALUES ($1, $2) 
ON CONFLICT (name) DO UPDATE SET title=$2, updated_at=now()`
	if _, err := r.db.ExecContext(ctx, q, name, sqlTitle); err != nil {
		return fmt.Errorf("db query failed: %w", err)
	}

	return nil
}
