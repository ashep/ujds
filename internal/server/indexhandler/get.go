package indexhandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"

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
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("index repo get failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&proto.GetResponse{
		Name:      index.Name,
		Title:     index.Title.String,
		Schema:    string(index.Schema),
		CreatedAt: uint64(index.CreatedAt.Unix()),
		UpdatedAt: uint64(index.UpdatedAt.Unix()),
	}), nil
}
