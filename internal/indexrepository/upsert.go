package indexrepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ashep/go-apperrors"
)

func (r *Repository) Upsert(ctx context.Context, name, schema string) error {
	if !r.nameRe.MatchString(name) {
		return apperrors.InvalidArgError{Subj: "name", Reason: "must match the regexp " + r.nameRe.String()}
	}

	nameParts := strings.Split(name, "/")
	for i := 1; i < len(nameParts); i++ {
		parentName := strings.Join(nameParts[:i], "/")

		_, err := r.Get(ctx, parentName)
		if errors.As(err, &apperrors.NotFoundError{}) {
			return apperrors.NotFoundError{Subj: "parent index " + parentName}
		} else if err != nil {
			return err
		}
	}

	if schema == "" {
		schema = "{}"
	}

	if err := json.Unmarshal([]byte(schema), &struct{}{}); err != nil {
		return apperrors.InvalidArgError{Subj: "schema", Reason: err.Error()}
	}

	q := `INSERT INTO index (name, schema) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET schema=$2, updated_at=now()`
	if _, err := r.db.ExecContext(ctx, q, name, schema); err != nil {
		return fmt.Errorf("db query failed: %w", err)
	}

	return nil
}
