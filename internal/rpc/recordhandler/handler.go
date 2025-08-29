package recordhandler

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
)

const perPageMax = 500

type indexRepo interface {
	Get(ctx context.Context, name string) (model.Index, error)
}

type recordRepo interface {
	Push(ctx context.Context, records []model.RecordUpdate) error
	Get(ctx context.Context, index string, id string) (model.Record, error)
	Find(ctx context.Context, index, search string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error)
	History(ctx context.Context, index, id string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error)
}

type Handler struct {
	ir  indexRepo
	rr  recordRepo
	now func() time.Time
	l   zerolog.Logger
}

func New(ir indexRepo, rr recordRepo, now func() time.Time, l zerolog.Logger) *Handler {
	return &Handler{ir: ir, rr: rr, now: now, l: l}
}
