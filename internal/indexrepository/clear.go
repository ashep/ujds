package indexrepository

import (
	"context"
	"fmt"
)

func (r *Repository) Clear(ctx context.Context, name string) error {
	if err := r.nameValidator.Validate(name); err != nil {
		return err //nolint:wrapcheck // ok
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`DELETE FROM record WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`, name)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("delete records: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`DELETE FROM record_log WHERE index_id=(SELECT id FROM index WHERE name=$1 LIMIT 1)`, name)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("delete record log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db commit: %w", err)
	}

	return nil
}
