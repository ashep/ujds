package handler

import (
	"context"

	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) CreateSchema(
	ctx context.Context,
	req *connect.Request[v1.CreateSchemaRequest],
) (*connect.Response[v1.CreateSchemaResponse], error) {
	ver, err := h.api.CreateSchema(ctx, req.Msg.Name, req.Msg.Schema)
	if err != nil {
		return nil, grpcErr(err, "api.CreateSchema failed", h.l.With().Str("proc", req.Spec().Procedure).Logger())
	}

	return connect.NewResponse(&v1.CreateSchemaResponse{Version: ver}), nil
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
		&v1.GetSchemaResponse{Name: sch.Name, Schema: string(sch.Schema)}), nil
}

func (h *Handler) UpdateSchema(
	ctx context.Context,
	req *connect.Request[v1.UpdateSchemaRequest],
) (*connect.Response[v1.UpdateSchemaResponse], error) {
	ver, err := h.api.UpdateSchema(ctx, req.Msg.Name, req.Msg.Schema, req.Msg.Version)
	if err != nil {
		return nil, grpcErr(err, "api.UpdateSchema failed", h.l.With().Str("proc", req.Spec().Procedure).Logger())
	}

	return connect.NewResponse(&v1.UpdateSchemaResponse{Version: ver}), nil
}
