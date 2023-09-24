package indexhandler

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"

	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func (h *Handler) List(
	ctx context.Context,
	req *connect.Request[proto.ListRequest],
) (*connect.Response[proto.ListResponse], error) {
	indices, err := h.repo.List(ctx)
	if err != nil {
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("index repo list failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	respData := make([]*proto.ListResponse_Index, 0)
	for _, idx := range indices {
		respData = append(respData, &proto.ListResponse_Index{Name: idx.Name, Title: idx.Title.String})
	}

	return connect.NewResponse(&proto.ListResponse{
		Indices: respData,
	}), nil
}
