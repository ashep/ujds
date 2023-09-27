package indexrepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ashep/go-apperrors"
)

func (r *Repository) Upsert(ctx context.Context, name, title, schema string) error {
	if err := r.nameValidator.Validate(name); err != nil {
		return apperrors.InvalidArgError{Subj: "name", Reason: err.Error()}
	}

	if schema == "" {
		schema = "{}"
	}

	if err := json.Unmarshal([]byte(schema), &struct{}{}); err != nil {
		return apperrors.InvalidArgError{Subj: "schema", Reason: err.Error()}
	}

	sqlTitle := sql.NullString{
		String: title,
		Valid:  title != "",
	}

	q := `INSERT INTO index (name, title, schema) VALUES ($1, $2, $3) 
ON CONFLICT (name) DO UPDATE SET title=$2, schema=$3, updated_at=now()`
	if _, err := r.db.ExecContext(ctx, q, name, sqlTitle, schema); err != nil {
		return fmt.Errorf("db query failed: %w", err)
	}

	return nil
}
