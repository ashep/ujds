package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"

	ujdsproto "github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func (h *Handler) PushIndex(
	ctx context.Context,
	req *connect.Request[ujdsproto.PushIndexRequest],
) (*connect.Response[ujdsproto.PushIndexResponse], error) {
	err := h.ir.Upsert(ctx, req.Msg.Name, req.Msg.Schema)
	if errors.As(err, &apperrors.InvalidArgError{}) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	} else if err != nil {
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("index repo upsert failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&ujdsproto.PushIndexResponse{}), nil
}

func (h *Handler) GetIndex(
	ctx context.Context,
	req *connect.Request[ujdsproto.GetIndexRequest],
) (*connect.Response[ujdsproto.GetIndexResponse], error) {
	index, err := h.ir.Get(ctx, req.Msg.Name)

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	case errors.As(err, &apperrors.NotFoundError{}):
		return nil, connect.NewError(connect.CodeNotFound, err)
	case err != nil:
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("index repo get failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&ujdsproto.GetIndexResponse{
		Name:      index.Name,
		Schema:    string(index.Schema),
		CreatedAt: uint64(index.CreatedAt.Unix()),
		UpdatedAt: uint64(index.UpdatedAt.Unix()),
	}), nil
}
