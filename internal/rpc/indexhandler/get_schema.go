package indexhandler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"

	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func (h *Handler) GetSchema(
	_ context.Context,
	req *connect.Request[proto.GetSchemaRequest],
) (*connect.Response[proto.GetSchemaResponse], error) {
	if err := h.nameValid.Validate(req.Msg.Name); err != nil {
		if errors.As(err, &apperrors.InvalidArgError{}) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		return nil, h.newInternalError(req, err, "index name validation failed")
	}

	schemas := h.schemas.SchemasFor(req.Msg.Name)

	res := &proto.GetSchemaResponse{
		Schemas: make([]*proto.GetSchemaResponse_Schema, 0, len(schemas)),
	}
	for _, s := range schemas {
		res.Schemas = append(res.Schemas, &proto.GetSchemaResponse_Schema{
			Pattern: s.Pattern,
			Schema:  string(s.Schema),
		})
	}

	return connect.NewResponse(res), nil
}
