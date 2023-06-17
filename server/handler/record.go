package handler

import (
	"context"
	"errors"
	"time"

	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushRecords(
	ctx context.Context,
	req *connect.Request[v1.PushRecordsRequest],
) (*connect.Response[v1.PushRecordsResponse], error) {
	if len(req.Msg.Records) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("empty records"))
	}

	apiRecords := make([]api.Record, 0)
	for _, rec := range req.Msg.Records {
		apiRecords = append(apiRecords, api.Record{
			Id:     rec.Id,
			Schema: rec.Schema,
			Data:   rec.Data,
		})
	}

	if err := h.api.PushRecords(ctx, apiRecords); err != nil {
		return nil, grpcErr(err, "api.PushRecords failed", h.l.With().Str("proc", req.Spec().Procedure).Logger())
	}

	return connect.NewResponse(&v1.PushRecordsResponse{}), nil
}

func (h *Handler) GetRecords(
	ctx context.Context,
	req *connect.Request[v1.GetRecordsRequest],
) (*connect.Response[v1.GetRecordsResponse], error) {
	if req.Msg.Schema == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("empty schema"))
	}

	since := time.Unix(req.Msg.Since, 0)
	records, cur, err := h.api.GetRecords(ctx, req.Msg.Schema, since, req.Msg.Cursor, req.Msg.Limit)
	if err != nil {
		return nil, grpcErr(err, "api.GetRecords failed", h.l.With().Str("proc", req.Spec().Procedure).Logger())
	}

	if len(records) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no records found"))
	}

	itemsR := make([]*v1.GetRecordsResponse_Record, len(records))
	for i, rec := range records {
		itemsR[i] = &v1.GetRecordsResponse_Record{
			Id:        rec.Id,
			Schema:    rec.Schema,
			Data:      rec.Data,
			CreatedAt: rec.CreatedAt.Unix(),
			UpdatedAt: rec.UpdatedAt.Unix(),
		}
	}

	return connect.NewResponse(&v1.GetRecordsResponse{Cursor: cur, Records: itemsR}), nil
}
