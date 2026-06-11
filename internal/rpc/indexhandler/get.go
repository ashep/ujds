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
// excluded from Get responses.
const catchAllPattern = ".*"

func (h *Handler) Get(
	ctx context.Context,
	req *connect.Request[proto.GetRequest],
) (*connect.Response[proto.GetResponse], error) {
	index, err := h.repo.Get(ctx, req.Msg.Name)

	switch {
	case errors.As(err, &apperrors.InvalidArgError{}):
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	case errors.As(err, &apperrors.NotFoundError{}):
		return nil, connect.NewError(connect.CodeNotFound, err)
	case err != nil:
		return nil, h.newInternalError(req, err, "index repo get failed")
	}

	schemas := h.schemas.SchemasFor(index.Name)
	res := &proto.GetResponse{
		Name:      index.Name,
		Title:     index.Title.String,
		CreatedAt: uint64(index.CreatedAt.Unix()), //nolint:gosec // ok
		UpdatedAt: uint64(index.UpdatedAt.Unix()), //nolint:gosec // ok
		Schemas:   make([]string, 0, len(schemas)),
	}
	for _, s := range schemas {
		if s.Pattern == catchAllPattern {
			continue
		}

		res.Schemas = append(res.Schemas, withDialect(s.Schema))
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
