package recordhandler

import (
	"context"
	"time"

	"github.com/ashep/ujds/internal/indexrepo"
	"github.com/ashep/ujds/internal/recordrepo"
	"github.com/rs/zerolog"
)

const perPageMax = 500

type indexRepo interface {
	Get(ctx context.Context, name string) (indexrepo.Index, error)
}

type recordRepo interface {
	Push(ctx context.Context, records []recordrepo.RecordUpdate) error
	Get(ctx context.Context, index string, id string) (recordrepo.Record, error)
	Find(ctx context.Context, req recordrepo.FindRequest) ([]recordrepo.Record, uint64, error)
	History(ctx context.Context, index, id string, since time.Time, cursor uint64, limit uint32) ([]recordrepo.Record, uint64, error)
}

type stringValidator interface {
	Validate(s string) error
}

type keyStringValidator interface {
	Validate(k, v string) error
}

type Handler struct {
	ir               indexRepo
	rr               recordRepo
	idxNameValidator stringValidator
	recIDValidator   stringValidator
	recJSONValidator keyStringValidator
	now              func() time.Time
	l                zerolog.Logger
}

func New(
	ir indexRepo,
	rr recordRepo,
	idxNameValidator,
	recIDValidator stringValidator,
	recDataValidator keyStringValidator,
	now func() time.Time,
	l zerolog.Logger,
) *Handler {
	return &Handler{
		ir:               ir,
		rr:               rr,
		idxNameValidator: idxNameValidator,
		recIDValidator:   recIDValidator,
		recJSONValidator: recDataValidator,
		now:              now,
		l:                l,
	}
}
