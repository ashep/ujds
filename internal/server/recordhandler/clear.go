package recordhandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"

	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func (h *Handler) Clear(
	ctx context.Context,
	req *connect.Request[proto.ClearRequest],
) (*connect.Response[proto.ClearResponse], error) {
	if err := h.rr.Clear(ctx, req.Msg.Index); err != nil {
		return nil, h.errAsConnect(err, req.Spec().Procedure, "ir.ClearRecords failed")
	}

	return connect.NewResponse(&proto.ClearResponse{}), nil
}

func (h *Handler) errAsConnect(err error, proc, msg string) error {
	switch {
	case errors.As(err, &apperrors.NotFoundError{}):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.As(err, &apperrors.InvalidArgError{}):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.As(err, &apperrors.AlreadyExistsError{}):
		return connect.NewError(connect.CodeAlreadyExists, err)
	case errors.As(err, &apperrors.AccessDeniedError{}):
		return connect.NewError(connect.CodeUnauthenticated, err)
	default:
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", proc).Int64("err_code", c).Msg(msg)

		return connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}
}
