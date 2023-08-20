package indexhandler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

//go:generate moq -out mock_test.go -pkg indexhandler_test -skip-ensure . indexRepo

type indexRepo interface {
	Upsert(ctx context.Context, name, schema string) error
	Get(ctx context.Context, name string) (model.Index, error)
}

type Handler struct {
	repo indexRepo
	now  func() time.Time
	l    zerolog.Logger
}

func New(repo indexRepo, now func() time.Time, l zerolog.Logger) *Handler {
	return &Handler{repo: repo, now: now, l: l}
}

func (h *Handler) Push(
	ctx context.Context,
	req *connect.Request[proto.PushRequest],
) (*connect.Response[proto.PushResponse], error) {
	err := h.repo.Upsert(ctx, req.Msg.Name, req.Msg.Schema)
	if errors.As(err, &apperrors.InvalidArgError{}) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	} else if err != nil {
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("index repo upsert failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&proto.PushResponse{}), nil
}

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
		c := h.now().Unix()
		h.l.Error().Err(err).Str("proc", req.Spec().Procedure).Int64("err_code", c).Msg("index repo get failed")

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("err_code: %d", c))
	}

	return connect.NewResponse(&proto.GetResponse{
		Name:      index.Name,
		Schema:    string(index.Schema),
		CreatedAt: uint64(index.CreatedAt.Unix()),
		UpdatedAt: uint64(index.UpdatedAt.Unix()),
	}), nil
}
