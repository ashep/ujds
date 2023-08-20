package recordhandler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

//go:generate moq -out mock_test.go -pkg recordhandler_test -skip-ensure . indexRepo recordRepo

type indexRepo interface {
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

func (h *Handler) Push(
	ctx context.Context,
	req *connect.Request[proto.PushRequest],
) (*connect.Response[proto.PushResponse], error) {
	index, err := h.ir.Get(ctx, req.Msg.Index)

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("index get failed: %w", err))
	case errors.As(err, &apperrors.NotFoundError{}):
		return nil, connect.NewError(connect.CodeNotFound, err)
	case err != nil:
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("index repo get failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	apiRecords := make([]model.Record, 0)
	for _, rec := range req.Msg.Records {
		apiRecords = append(apiRecords, model.Record{
			ID:   rec.Id,
			Data: rec.Data,
		})
	}

	err = h.rr.Push(ctx, index, apiRecords)
	if errors.As(err, &apperrors.InvalidArgError{}) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	} else if err != nil {
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("record repo push failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&proto.PushResponse{}), nil
}

func (h *Handler) Get(
	ctx context.Context,
	req *connect.Request[proto.GetRequest],
) (*connect.Response[proto.GetResponse], error) {
	rec, err := h.rr.Get(ctx, req.Msg.Index, req.Msg.Id)

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	case errors.As(err, &apperrors.NotFoundError{}):
		return nil, connect.NewError(connect.CodeNotFound, err)
	case err != nil:
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("record repo push failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&proto.GetResponse{Record: &proto.Record{
		Id:        rec.ID,
		Rev:       rec.Rev,
		Index:     rec.Index,
		CreatedAt: rec.CreatedAt.Unix(),
		UpdatedAt: rec.UpdatedAt.Unix(),
		Data:      rec.Data,
	}}), nil
}

func (h *Handler) GetAll(
	ctx context.Context,
	req *connect.Request[proto.GetAllRequest],
) (*connect.Response[proto.GetAllResponse], error) {
	if req.Msg.Index == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("index is not specified"))
	}

	since := time.Unix(req.Msg.Since, 0)

	records, cur, err := h.rr.GetAll(ctx, req.Msg.Index, since, req.Msg.Cursor, req.Msg.Limit)
	if err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "ir.GetAll failed")
	}

	if len(records) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no records found"))
	}

	itemsR := make([]*proto.Record, len(records))
	for i, rec := range records {
		itemsR[i] = &proto.Record{
			Id:        rec.ID,
			Rev:       rec.Rev,
			Index:     rec.Index,
			Data:      rec.Data,
			CreatedAt: rec.CreatedAt.Unix(),
			UpdatedAt: rec.UpdatedAt.Unix(),
		}
	}

	return connect.NewResponse(&proto.GetAllResponse{Cursor: cur, Records: itemsR}), nil
}

func (h *Handler) Clear(
	ctx context.Context,
	req *connect.Request[proto.ClearRequest],
) (*connect.Response[proto.ClearResponse], error) {
	if err := h.rr.Clear(ctx, req.Msg.Index); err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "ir.ClearRecords failed")
	}

	return connect.NewResponse(&proto.ClearResponse{}), nil
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
