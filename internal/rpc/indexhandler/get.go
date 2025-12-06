package indexhandler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"

	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func (h *Handler) Get(
	ctx context.Context,
	req *connect.Request[proto.GetRequest],
) (*connect.Response[proto.GetResponse], error) {
	index, err := h.repo.Get(ctx, req.Msg.Name)

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	case errors.As(err, &apperrors.NotFoundError{}):
		return nil, connect.NewError(connect.CodeNotFound, err)
	case err != nil:
		return nil, h.newInternalError(req, err, "index repo get failed")
	}

	return connect.NewResponse(&proto.GetResponse{
		Name:      index.Name,
		Title:     index.Title.String,
		CreatedAt: uint64(index.CreatedAt.Unix()), //nolint:gosec // ok
		UpdatedAt: uint64(index.UpdatedAt.Unix()), //nolint:gosec // ok
	}), nil
}
