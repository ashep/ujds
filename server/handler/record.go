package handler

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/errs"
	"github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushRecords(
	ctx context.Context,
	req *connect.Request[v1.PushRecordsRequest],
) (*connect.Response[v1.PushRecordsResponse], error) {
	if len(req.Msg.Records) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("no records provided"))
	}

	apiRecords := make([]api.Record, 0)
	for _, rec := range req.Msg.Records {
		apiRecords = append(apiRecords, api.Record{
			ID:     rec.Id,
			Schema: rec.Schema,
			Data:   rec.Data,
		})
	}

	if err := h.api.PushRecords(ctx, apiRecords); errors.Is(err, errs.ErrEmptyArg{}) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
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

	records, cur, err := h.api.GetRecords(ctx, req.Msg.Schema, req.Msg.Cursor, req.Msg.Limit)
	if errors.Is(err, errs.ErrNotFound{}) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not found"))
	} else if err != nil {
		h.l.Error().
			Err(err).
			Str("schema", req.Msg.Schema).
			Msg("get records")
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	itemsR := make([]*v1.GetRecordsResponse_Record, len(records))
	for i, item := range records {
		itemsR[i] = &v1.GetRecordsResponse_Record{
			Id:      item.ID,
			Version: item.Version,
			Data:    item.Data,
		}
	}

	return connect.NewResponse(&v1.GetRecordsResponse{Cursor: cur, Records: itemsR}), nil
}
