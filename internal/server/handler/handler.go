package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/indexrepository"
	"github.com/ashep/ujds/internal/recordrepository"
)

//go:generate moq -out handler_mock_test.go -pkg handler_test -skip-ensure . indexRepo recordRepo

type indexRepo interface {
	Upsert(ctx context.Context, name, schema string) error
	Get(ctx context.Context, name string) (indexrepository.Index, error)
}

type recordRepo interface {
	Push(ctx context.Context, index indexrepository.Index, records []recordrepository.Record) error
	Get(ctx context.Context, indexName string, id string) (recordrepository.Record, error)
	GetAll(ctx context.Context, indexName string, since time.Time, cursor uint64, limit uint32) ([]recordrepository.Record, uint64, error)
	Clear(ctx context.Context, indexName string) error
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
	case errors.As(err, &apperrors.InvalidArgError{}), errors.As(err, &apperrors.EmptyArgError{}):
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
