package recordhandler

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
)

//go:generate moq -out mock_test.go -pkg recordhandler_test -skip-ensure . indexRepo recordRepo

type indexRepo interface {
	Get(ctx context.Context, name string) (model.Index, error)
}

type recordRepo interface {
	Push(ctx context.Context, indexID uint, schema []byte, records []model.Record) error
	Get(ctx context.Context, index string, id string) (model.Record, error)
	GetAll(ctx context.Context, index string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error)
	Clear(ctx context.Context, index string) error
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
