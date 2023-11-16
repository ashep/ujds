package recordrepository

import (
	"context"
	"time"

	"github.com/ashep/ujds/internal/model"
)

func (r *Repository) Search(
	ctx context.Context,
	query string,
	index string,
	since time.Time,
	cursor uint64,
	limit uint32,
) ([]model.Record, uint64, error) {
	// q := `SELECT data FROM record WHERE data->>`

	return nil, 0, nil
}
