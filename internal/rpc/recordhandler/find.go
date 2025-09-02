package recordhandler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"
	"github.com/ashep/ujds/internal/recordrepo"

	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func (h *Handler) Find(
	ctx context.Context,
	req *connect.Request[proto.FindRequest],
) (*connect.Response[proto.FindResponse], error) {
	if req.Msg.Limit == 0 || req.Msg.Limit > perPageMax {
		req.Msg.Limit = perPageMax
	}

	var ntSince *time.Time
	if req.Msg.NotTouchedSince != 0 {
		t := time.Unix(req.Msg.NotTouchedSince, 0)
		ntSince = &t
	}

	records, cur, err := h.rr.Find(ctx, recordrepo.FindRequest{
		Index:           req.Msg.Index,
		Query:           req.Msg.Search,
		Since:           time.Unix(req.Msg.Since, 0),
		Cursor:          req.Msg.Cursor,
		Limit:           req.Msg.Limit,
		NotTouchedSince: ntSince,
	})

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	case err != nil:
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("record repo find failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	itemsR := make([]*proto.Record, len(records))
	for i, rec := range records {
		itemsR[i] = &proto.Record{
			Id:        rec.ID,
			Rev:       rec.Rev,
			Index:     req.Msg.Index,
			Data:      rec.Data,
			CreatedAt: rec.CreatedAt.Unix(),
			UpdatedAt: rec.UpdatedAt.Unix(),
			TouchedAt: rec.TouchedAt.Unix(),
		}
	}

	return connect.NewResponse(&proto.FindResponse{Cursor: cur, Records: itemsR}), nil
}
