package indexhandler

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"connectrpc.com/connect"

	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func (h *Handler) List(
	ctx context.Context,
	req *connect.Request[proto.ListRequest],
) (*connect.Response[proto.ListResponse], error) {
	patterns, err := h.makeIndexNameFilter(req.Msg.GetFilter())
	if err != nil {
		return nil, h.newInternalError(req, err, "index name filter build failed")
	}

	indices, err := h.repo.List(ctx)
	if err != nil {
		return nil, h.newInternalError(req, err, "index repo list failed")
	}

	respData := make([]*proto.ListResponse_Index, 0)

	for _, idx := range indices {
		if h.filterIndexName(patterns, idx.Name) {
			respData = append(respData, &proto.ListResponse_Index{Name: idx.Name, Title: idx.Title.String})
		}
	}

	return connect.NewResponse(&proto.ListResponse{
		Indices: respData,
	}), nil
}

func (h *Handler) makeIndexNameFilter(reqFilter *proto.ListRequestFilter) ([]*regexp.Regexp, error) {
	if reqFilter == nil {
		return nil, nil
	}

	names := reqFilter.GetNames()
	if len(names) == 0 {
		return nil, nil
	}

	patterns := make([]*regexp.Regexp, len(names))

	for i, pat := range names {
		pat = strings.ReplaceAll(pat, ".", "\\.")
		pat = strings.ReplaceAll(pat, "*", ".*")

		re, err := regexp.Compile("^" + pat + "$")
		if err != nil {
			return nil, fmt.Errorf("invalid index name pattern: %w", err)
		}

		patterns[i] = re
	}

	return patterns, nil
}

func (h *Handler) filterIndexName(patterns []*regexp.Regexp, name string) bool {
	if len(patterns) == 0 {
		return true
	}

	for _, p := range patterns {
		if p.MatchString(name) {
			return true
		}
	}

	return false
}
