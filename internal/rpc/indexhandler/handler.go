package indexhandler

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
)

type indexRepo interface {
	Upsert(ctx context.Context, name, title, schema string) error
	Get(ctx context.Context, name string) (model.Index, error)
	List(ctx context.Context) ([]model.Index, error)
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

func (h *Handler) newInternalError(req connect.AnyRequest, err error, msg string) error {
	c := h.now().Unix()
	h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg(msg)

	return connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
}
