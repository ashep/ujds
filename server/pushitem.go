package server

import (
	"context"
	"errors"

	"github.com/ashep/datapimp/dataservice"
	v1 "github.com/ashep/datapimp/gen/proto/datapimp/v1"
	"github.com/bufbuild/connect-go"
)

func (s *Server) PushItem(
	ctx context.Context,
	req *connect.Request[v1.PushItemRequest],
) (*connect.Response[v1.PushItemResponse], error) {
	s.l.Debug().
		Str("type", req.Msg.Type).
		Str("id", req.Msg.Id).
		Msg("push data item request")

	if req.Msg.Type == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("type is empty"))
	}

	i, err := s.ds.Push(ctx, req.Msg.Type, req.Msg.Id, []byte(req.Msg.Data))
	if errors.Is(err, dataservice.ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not found"))
	} else if err != nil {
		s.l.Error().
			Err(err).
			Str("type", req.Msg.Type).
			Str("id", req.Msg.Id).
			Str("data", req.Msg.Data).
			Msg("failed to push an item")
		return nil, connect.NewError(connect.CodeInternal, errors.New("operation failed"))
	}

	item := &v1.Item{
		Id:      i.Id,
		Type:    i.Type,
		Version: i.Version,
		Time:    uint64(i.Time.Unix()),
		Data:    string(i.Data),
	}

	return connect.NewResponse(&v1.PushItemResponse{
		Item: item,
	}), nil
}
