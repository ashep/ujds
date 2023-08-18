package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
)

//go:generate moq -out handler_mock_test.go -pkg handler_test -skip-ensure . indexRepo recordRepo

type indexRepo interface {
	Upsert(ctx context.Context, name, schema string) error
	Get(ctx context.Context, name string) (model.Index, error)
}

type recordRepo interface {
	Push(ctx context.Context, index model.Index, records []model.Record) error
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

func (h *Handler) errAsConnect(err error, proc, msg string) error {
	switch {
	case errors.As(err, &apperrors.NotFoundError{}):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.As(err, &apperrors.InvalidArgError{}):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.As(err, &apperrors.AlreadyExistsError{}):
		return connect.NewError(connect.CodeAlreadyExists, err)
	case errors.As(err, &apperrors.AccessDeniedError{}):
		return connect.NewError(connect.CodeUnauthenticated, err)
	default:
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", proc).Int64("err_code", c).Msg(msg)

		return connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}
}
