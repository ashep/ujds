package indexhandler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"

	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func (h *Handler) Clear(
	ctx context.Context,
	req *connect.Request[proto.ClearRequest],
) (*connect.Response[proto.ClearResponse], error) {
	err := h.repo.Clear(ctx, req.Msg.Name)

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	case err != nil:
		return nil, h.newInternalError(req, err, "index repo clear failed")
	}

	return connect.NewResponse(&proto.ClearResponse{}), nil
}
