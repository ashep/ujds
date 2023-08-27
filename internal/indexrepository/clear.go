package indexrepository

import (
	"context"
	"fmt"

	"github.com/ashep/go-apperrors"
)

func (r *Repository) Clear(ctx context.Context, name string) error {
	if !r.nameRe.MatchString(name) {
		return apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp " + r.nameRe.String()}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM record WHERE index_id=$1`, name)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete records: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM record_log WHERE index_id=$1`, name)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete record logs: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db commit failed: %w", err)
	}

	return nil
}
