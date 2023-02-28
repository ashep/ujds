package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ashep/datapimp/dataservice"
	"github.com/ashep/datapimp/gen/proto/datapimp/v1"
	"github.com/bufbuild/connect-go"
)

func (h *Handler) GetItem(
	ctx context.Context,
	req *connect.Request[v1.GetItemRequest],
) (*connect.Response[v1.GetItemResponse], error) {
	h.l.Debug().Str("id", req.Msg.Id).Msg("get data item request")

	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is empty"))
	}

	i, err := h.data.GetItem(ctx, req.Msg.Id)
	if errors.Is(err, dataservice.ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not found"))
	} else if err != nil {
		h.l.Error().
			Err(err).
			Str("id", req.Msg.Id).
			Msg("failed to get a data item")
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	item := &v1.Item{
		Id:      i.Id,
		Type:    i.Type,
		Version: i.Version,
		Time:    uint64(i.Time.Unix()),
		Data:    string(i.Data),
	}

	return connect.NewResponse(&v1.GetItemResponse{
		Item: item,
	}), nil
}

func (h *Handler) PushItem(
	ctx context.Context,
	req *connect.Request[v1.PushItemRequest],
) (*connect.Response[v1.PushItemResponse], error) {
	h.l.Debug().
		Str("type", req.Msg.Type).
		Str("id", req.Msg.Id).
		Msg("push data item request")

	if req.Msg.Type == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("type is empty"))
	}

	if req.Msg.Data == "" {
		req.Msg.Data = "{}"
	}

	d := make(map[string]any)
	err := json.Unmarshal([]byte(req.Msg.Data), &d)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to unmarshal data: %w", err))
	}

	i, err := h.data.UpsertItem(ctx, req.Msg.Type, req.Msg.Id, []byte(req.Msg.Data))
	if errors.Is(err, dataservice.ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not found"))
	} else if err != nil {
		h.l.Error().
			Err(err).
			Str("type", req.Msg.Type).
			Str("id", req.Msg.Id).
			Str("data", req.Msg.Data).
			Msg("failed to push an item")
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	item := &v1.Item{
		Id:      i.Id,
		Type:    i.Type,
		Version: i.Version,
		Time:    uint64(i.Time.Unix()),
		Data:    string(i.Data),
	}

	return connect.NewResponse(&v1.PushItemResponse{
		Item: item,
	}), nil
}
