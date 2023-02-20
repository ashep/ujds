package server

import (
	"context"
	"errors"

	"github.com/ashep/datapimp/dataservice"
	v1 "github.com/ashep/datapimp/gen/proto/datapimp/v1"
	"github.com/bufbuild/connect-go"
)

func (s *Server) GetItem(
	ctx context.Context,
	req *connect.Request[v1.GetItemRequest],
) (*connect.Response[v1.GetItemResponse], error) {
	s.l.Debug().Str("id", req.Msg.Id).Msg("get data item request")

	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is empty"))
	}

	i, err := s.ds.GetItem(ctx, req.Msg.Id)
	if errors.Is(err, dataservice.ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not found"))
	} else if err != nil {
		s.l.Error().
			Err(err).
			Str("id", req.Msg.Id).
			Msg("failed to get a data item")
		return nil, connect.NewError(connect.CodeInternal, errors.New("operation failed"))
	}

	item := &v1.Item{
		Id:      i.Id,
		Type:    i.Type,
		Version: i.Version,
		Time:    uint64(i.Time.Unix()),
		Data:    string(i.Data),
	}

	return connect.NewResponse(&v1.GetItemResponse{
		Item: item,
	}), nil
}
