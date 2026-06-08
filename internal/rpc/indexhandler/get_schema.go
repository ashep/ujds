package indexhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"

	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

// draft07Schema is the dialect the server validates against (xeipuuv/gojsonschema
// supports up to draft-07). It is stamped onto returned schemas that don't declare one.
const draft07Schema = "http://json-schema.org/draft-07/schema#"

// catchAllPattern is the synthetic ".*" -> "{}" entry the config always injects so that
// every record is at least valid JSON. It carries no validation constraints, so it is
// excluded from GetSchema responses.
const catchAllPattern = ".*"

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
		if s.Pattern == catchAllPattern {
			continue
		}

		res.Schemas = append(res.Schemas, &proto.GetSchemaResponse_Schema{
			Pattern: s.Pattern,
			Schema:  withDialect(s.Schema),
		})
	}

	return connect.NewResponse(res), nil
}

// withDialect returns the schema with a "$schema" dialect declaration. If the schema
// already declares one it is returned unchanged; otherwise the draft-07 dialect is
// prepended (the original key order of the rest is preserved). Non-object schemas
// (e.g. malformed input or a boolean schema) are returned as-is.
func withDialect(raw json.RawMessage) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil || obj == nil {
		return string(raw)
	}

	if _, ok := obj["$schema"]; ok {
		return string(raw)
	}

	if len(obj) == 0 {
		return `{"$schema":"` + draft07Schema + `"}`
	}

	inner := bytes.TrimSpace(raw)

	return `{"$schema":"` + draft07Schema + `",` + string(inner[1:])
}
