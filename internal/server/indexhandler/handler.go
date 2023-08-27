package indexhandler

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
)

//go:generate moq -out mock_test.go -pkg indexhandler_test -skip-ensure . indexRepo

type indexRepo interface {
	Upsert(ctx context.Context, name, schema string) error
	Get(ctx context.Context, name string) (model.Index, error)
	Clear(ctx context.Context, name string) error
}

type Handler struct {
	repo indexRepo
	now  func() time.Time
	l    zerolog.Logger
}

func New(repo indexRepo, now func() time.Time, l zerolog.Logger) *Handler {
	return &Handler{repo: repo, now: now, l: l}
}
