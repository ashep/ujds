package handler

import (
	"context"

	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushSchema(
	ctx context.Context,
	req *connect.Request[v1.PushSchemaRequest],
) (*connect.Response[v1.PushSchemaResponse], error) {
	if err := h.api.PushSchema(ctx, req.Msg.Name, req.Msg.Data); err != nil {
		return nil, grpcErr(err, "api.PushSchema failed", h.l.With().Str("proc", req.Spec().Procedure).Logger())
	}

	return connect.NewResponse(&v1.PushSchemaResponse{}), nil
}

func (h *Handler) GetSchema(
	ctx context.Context,
	req *connect.Request[v1.GetSchemaRequest],
) (*connect.Response[v1.GetSchemaResponse], error) {
	var sch api.Schema
	var err error

	sch, err = h.api.GetSchema(ctx, req.Msg.Name)

	if err != nil {
		return nil, grpcErr(err, "api.GetSchema failed", h.l.With().Str("proc", req.Spec().Procedure).Logger())
	}

	return connect.NewResponse(
		&v1.GetSchemaResponse{
			Name:      sch.Name,
			Data:      string(sch.Data),
			CreatedAt: uint64(sch.CreatedAt.Unix()),
		}), nil
}
