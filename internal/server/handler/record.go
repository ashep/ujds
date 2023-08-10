package handler

import (
	"context"
	"errors"
	"time"

	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/internal/recordrepository"
	v1 "github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushRecords(
	ctx context.Context,
	req *connect.Request[v1.PushRecordsRequest],
) (*connect.Response[v1.PushRecordsResponse], error) {
	if req.Msg.Index == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("index is not specified"))
	}

	index, err := h.ir.Get(ctx, req.Msg.Index)
	if err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "index get failed")
	}

	if len(req.Msg.Records) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("empty records"))
	}

	apiRecords := make([]recordrepository.Record, 0)
	for _, rec := range req.Msg.Records {
		apiRecords = append(apiRecords, recordrepository.Record{
			ID:   rec.Id,
			Data: rec.Data,
		})
	}

	if err := h.rr.Push(ctx, index, apiRecords); err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "ir.PushRecords failed")
	}

	return connect.NewResponse(&v1.PushRecordsResponse{}), nil
}

func (h *Handler) GetRecord(
	ctx context.Context,
	req *connect.Request[v1.GetRecordRequest],
) (*connect.Response[v1.GetRecordResponse], error) {
	if req.Msg.Index == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("index is not specified"))
	}

	rec, err := h.rr.Get(ctx, req.Msg.Index, req.Msg.Id)
	if err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "ir.ClearRecords failed")
	}

	return connect.NewResponse(&v1.GetRecordResponse{Record: &v1.Record{
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
	req *connect.Request[v1.GetRecordsRequest],
) (*connect.Response[v1.GetRecordsResponse], error) {
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

	itemsR := make([]*v1.Record, len(records))
	for i, rec := range records {
		itemsR[i] = &v1.Record{
			Id:        rec.ID,
			Rev:       rec.Rev,
			Index:     rec.Index,
			Data:      rec.Data,
			CreatedAt: rec.CreatedAt.Unix(),
			UpdatedAt: rec.UpdatedAt.Unix(),
		}
	}

	return connect.NewResponse(&v1.GetRecordsResponse{Cursor: cur, Records: itemsR}), nil
}

func (h *Handler) ClearRecords(
	ctx context.Context,
	req *connect.Request[v1.ClearRecordsRequest],
) (*connect.Response[v1.ClearRecordsResponse], error) {
	if err := h.rr.Clear(ctx, req.Msg.Index); err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "ir.ClearRecords failed")
	}

	return connect.NewResponse(&v1.ClearRecordsResponse{}), nil
}
