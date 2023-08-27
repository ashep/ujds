package recordhandler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"

	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func (h *Handler) GetAll(
	ctx context.Context,
	req *connect.Request[proto.GetAllRequest],
) (*connect.Response[proto.GetAllResponse], error) {
	since := time.Unix(req.Msg.Since, 0)
	records, cur, err := h.rr.GetAll(ctx, req.Msg.Index, since, req.Msg.Cursor, req.Msg.Limit)

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	case err != nil:
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("record repo get all failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	if len(records) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no records found"))
	}

	itemsR := make([]*proto.Record, len(records))
	for i, rec := range records {
		itemsR[i] = &proto.Record{
			Id:        rec.ID,
			Rev:       rec.Rev,
			Index:     rec.Index,
			Data:      rec.Data,
			CreatedAt: rec.CreatedAt.Unix(),
			UpdatedAt: rec.UpdatedAt.Unix(),
		}
	}

	return connect.NewResponse(&proto.GetAllResponse{Cursor: cur, Records: itemsR}), nil
}
