package recordrepository

import (
	"context"
	"fmt"
)

func (r *Repository) Clear(ctx context.Context, indexName string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM record WHERE index_id=$1`, indexName)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete records: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM record_log WHERE index_id=$1`, indexName)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete records: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db commit failed: %w", err)
	}

	return nil
}
