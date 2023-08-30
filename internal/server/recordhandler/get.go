package recordhandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"

	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func (h *Handler) Get(
	ctx context.Context,
	req *connect.Request[proto.GetRequest],
) (*connect.Response[proto.GetResponse], error) {
	rec, err := h.rr.Get(ctx, req.Msg.Index, req.Msg.Id)

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	case errors.As(err, &apperrors.NotFoundError{}):
		return nil, connect.NewError(connect.CodeNotFound, err)
	case err != nil:
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("record repo push failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&proto.GetResponse{Record: &proto.Record{
		Id:        rec.ID,
		Rev:       rec.Rev,
		Index:     req.Msg.Index,
		CreatedAt: rec.CreatedAt.Unix(),
		UpdatedAt: rec.UpdatedAt.Unix(),
		Data:      rec.Data,
	}}), nil
}
