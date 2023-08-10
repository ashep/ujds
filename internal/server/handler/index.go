package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"

	v1 "github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushIndex(
	ctx context.Context,
	req *connect.Request[v1.PushIndexRequest],
) (*connect.Response[v1.PushIndexResponse], error) {
	err := h.ir.Upsert(ctx, req.Msg.Name, req.Msg.Schema)
	if errors.As(err, &apperrors.InvalidArgError{}) || errors.As(err, &apperrors.EmptyArgError{}) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	} else if err != nil {
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("index repo upsert failed")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&v1.PushIndexResponse{}), nil
}

func (h *Handler) GetIndex(
	ctx context.Context,
	req *connect.Request[v1.GetIndexRequest],
) (*connect.Response[v1.GetIndexResponse], error) {
	index, err := h.ir.Get(ctx, req.Msg.Name)
	if err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "index repo get failed")
	}

	return connect.NewResponse(&v1.GetIndexResponse{
		Name:      index.Name,
		Schema:    string(index.Schema),
		CreatedAt: uint64(index.CreatedAt.Unix()),
		UpdatedAt: uint64(index.UpdatedAt.Unix()),
	}), nil
}
