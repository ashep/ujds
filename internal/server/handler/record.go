package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/internal/model"
	ujdsproto "github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushRecords(
	ctx context.Context,
	req *connect.Request[ujdsproto.PushRecordsRequest],
) (*connect.Response[ujdsproto.PushRecordsResponse], error) {
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

	return connect.NewResponse(&ujdsproto.PushRecordsResponse{}), nil
}

func (h *Handler) GetRecord(
	ctx context.Context,
	req *connect.Request[ujdsproto.GetRecordRequest],
) (*connect.Response[ujdsproto.GetRecordResponse], error) {
	if req.Msg.Index == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("index is not specified"))
	}

	rec, err := h.rr.Get(ctx, req.Msg.Index, req.Msg.Id)
	if err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "ir.ClearRecords failed")
	}

	return connect.NewResponse(&ujdsproto.GetRecordResponse{Record: &ujdsproto.Record{
		Id:        rec.ID,
		Rev:       rec.Rev,
		Index:     rec.Index,
		CreatedAt: rec.CreatedAt.Unix(),
		UpdatedAt: rec.UpdatedAt.Unix(),
		Data:      rec.Data,
	}}), nil
}

func (h *Handler) GetRecords(
	ctx context.Context,
	req *connect.Request[ujdsproto.GetRecordsRequest],
) (*connect.Response[ujdsproto.GetRecordsResponse], error) {
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

	itemsR := make([]*ujdsproto.Record, len(records))
	for i, rec := range records {
		itemsR[i] = &ujdsproto.Record{
			Id:        rec.ID,
			Rev:       rec.Rev,
			Index:     rec.Index,
			Data:      rec.Data,
			CreatedAt: rec.CreatedAt.Unix(),
			UpdatedAt: rec.UpdatedAt.Unix(),
		}
	}

	return connect.NewResponse(&ujdsproto.GetRecordsResponse{Cursor: cur, Records: itemsR}), nil
}

func (h *Handler) ClearRecords(
	ctx context.Context,
	req *connect.Request[ujdsproto.ClearRecordsRequest],
) (*connect.Response[ujdsproto.ClearRecordsResponse], error) {
	if err := h.rr.Clear(ctx, req.Msg.Index); err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "ir.ClearRecords failed")
	}

	return connect.NewResponse(&ujdsproto.ClearRecordsResponse{}), nil
}
