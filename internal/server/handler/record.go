package handler

import (
	"context"
	"errors"
	"time"

	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/internal/api"
	"github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushRecords(
	ctx context.Context,
	req *connect.Request[v1.PushRecordsRequest],
) (*connect.Response[v1.PushRecordsResponse], error) {
	if req.Msg.Index == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("index is not specified"))
	}

	if len(req.Msg.Records) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("empty records"))
	}

	apiRecords := make([]api.Record, 0)
	for _, rec := range req.Msg.Records {
		apiRecords = append(apiRecords, api.Record{
			Id:   rec.Id,
			Data: rec.Data,
		})
	}

	if err := h.api.PushRecords(ctx, req.Msg.Index, apiRecords); err != nil {
		return nil, grpcErr(err, req.Spec().Procedure, "api.PushRecords failed", h.l)
	}

	return connect.NewResponse(&v1.PushRecordsResponse{}), nil
}

func (h *Handler) GetRecords(
	ctx context.Context,
	req *connect.Request[v1.GetRecordsRequest],
) (*connect.Response[v1.GetRecordsResponse], error) {
	if req.Msg.Index == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("index is not specified"))
	}

	since := time.Unix(req.Msg.Since, 0)
	records, cur, err := h.api.GetRecords(ctx, req.Msg.Index, since, req.Msg.Cursor, req.Msg.Limit)
	if err != nil {
		return nil, grpcErr(err, req.Spec().Procedure, "api.GetRecords failed", h.l)
	}

	if len(records) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no records found"))
	}

	itemsR := make([]*v1.GetRecordsResponse_Record, len(records))
	for i, rec := range records {
		itemsR[i] = &v1.GetRecordsResponse_Record{
			Id:        rec.Id,
			Rev:       rec.Rev,
			Index:     rec.Index,
			Data:      rec.Data,
			CreatedAt: rec.CreatedAt.Unix(),
			UpdatedAt: rec.UpdatedAt.Unix(),
		}
	}

	return connect.NewResponse(&v1.GetRecordsResponse{Cursor: cur, Records: itemsR}), nil
}
