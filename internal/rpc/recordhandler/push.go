package recordhandler

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"
	"github.com/ashep/ujds/internal/indexrepo"
	"github.com/ashep/ujds/internal/recordrepo"

	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func (h *Handler) Push(
	ctx context.Context,
	req *connect.Request[proto.PushRequest],
) (*connect.Response[proto.PushResponse], error) {
	if len(req.Msg.GetRecords()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("empty records"))
	}

	cache := make(map[string]indexrepo.Index)
	updates := make([]recordrepo.RecordUpdate, 0)

	for _, rec := range req.Msg.GetRecords() {
		index, err := h.getIndex(ctx, req.Spec().Procedure, rec.Index, cache)
		if err != nil {
			return nil, err
		}

		updates = append(updates, recordrepo.RecordUpdate{
			ID:      rec.GetId(),
			IndexID: index.ID,
			Schema:  index.Schema,
			Data:    rec.GetData(),
		})
	}

	err := h.rr.Push(ctx, updates)
	if errors.As(err, &apperrors.InvalidArgError{}) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	} else if err != nil {
		c := h.now().UnixMilli()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("record repo push failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&proto.PushResponse{}), nil
}

func (h *Handler) getIndex(ctx context.Context, proc, name string, cache map[string]indexrepo.Index) (indexrepo.Index, error) {
	var err error

	index, ok := cache[name]

	if ok {
		return index, nil
	}

	if index, err = h.ir.Get(ctx, name); err != nil {
		switch {
		case errors.As(err, &apperrors.InvalidArgError{}):
			return indexrepo.Index{}, connect.NewError(connect.CodeInvalidArgument, err)
		case errors.As(err, &apperrors.NotFoundError{}):
			return indexrepo.Index{}, connect.NewError(connect.CodeNotFound, err)
		default:
			c := h.now().UnixMilli()
			h.l.Error().Err(err).Str("proc", proc).Int64("err_code", c).Msg("index repo get failed")

			return indexrepo.Index{}, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
		}
	}

	return index, nil
}
