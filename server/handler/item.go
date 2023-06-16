package handler

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/errs"
	"github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) SetItems(
	ctx context.Context,
	req *connect.Request[v1.SetItemsRequest],
) (*connect.Response[v1.SetItemsResponse], error) {
	if len(req.Msg.Items) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("empty items"))
	}

	apiItems := make([]api.Item, 0)
	for _, item := range req.Msg.Items {
		apiItems = append(apiItems, api.Item{
			ID:     item.Id,
			Schema: item.Schema,
			Data:   item.Data,
		})
	}

	if err := h.api.InsertItems(ctx, apiItems); errors.Is(err, errs.ErrEmptyArg{}) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&v1.SetItemsResponse{}), nil
}

func (h *Handler) GetItems(
	ctx context.Context,
	req *connect.Request[v1.GetItemsRequest],
) (*connect.Response[v1.GetItemsResponse], error) {
	if req.Msg.Schema == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("empty schema"))
	}

	items, cur, err := h.api.GetRecords(ctx, req.Msg.Schema, req.Msg.Version, req.Msg.Cursor, req.Msg.Limit)
	if errors.Is(err, errs.ErrNotFound{}) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not found"))
	} else if err != nil {
		h.l.Error().
			Err(err).
			Str("schema", req.Msg.Schema).
			Uint64("version", req.Msg.Version).
			Msg("get items")
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	itemsR := make([]*v1.GetItemsResponse_Item, len(items))
	for i, item := range items {
		itemsR[i] = &v1.GetItemsResponse_Item{
			Id:      item.ID,
			Version: item.Version,
			Data:    item.Data,
		}
	}

	return connect.NewResponse(&v1.GetItemsResponse{
		Cursor: cur,
		Items:  itemsR,
	}), nil
}
