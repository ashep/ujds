package handler

import (
	"context"

	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushIndex(
	ctx context.Context,
	req *connect.Request[v1.PushIndexRequest],
) (*connect.Response[v1.PushIndexResponse], error) {
	if err := h.api.UpsertIndex(ctx, req.Msg.Name, req.Msg.Schema); err != nil {
		return nil, grpcErr(err, req.Spec().Procedure, "api.PushIndex failed", h.l)
	}

	return connect.NewResponse(&v1.PushIndexResponse{}), nil
}

func (h *Handler) GetIndex(
	ctx context.Context,
	req *connect.Request[v1.GetIndexRequest],
) (*connect.Response[v1.GetIndexResponse], error) {
	index, err := h.api.GetIndex(ctx, req.Msg.Name)
	if err != nil {
		return nil, grpcErr(err, req.Spec().Procedure, "api.GetIndex failed", h.l)
	}

	return connect.NewResponse(&v1.GetIndexResponse{
		Name:      index.Name,
		Schema:    string(index.Data),
		CreatedAt: uint64(index.CreatedAt.Unix()),
		UpdatedAt: uint64(index.UpdatedAt.Unix()),
	}), nil
}
