package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ashep/datapimp/authservice"
	"github.com/ashep/datapimp/gen/proto/datapimp/v1"
	"github.com/bufbuild/connect-go"
)

func (h *Handler) CreateEntity(
	ctx context.Context,
	req *connect.Request[v1.CreateEntityRequest],
) (*connect.Response[v1.CreateEntityResponse], error) {
	authId := ctx.Value("authId")
	if authId != "admin" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	authSec, ok := ctx.Value("authSecret").(string)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	entitySec := req.Msg.Secret
	if len(entitySec) < 8 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("secret is too short"))
	}

	if req.Msg.Permissions == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("empty permissions"))
	}

	prm := authservice.Permissions{}
	if err := json.Unmarshal([]byte(req.Msg.Permissions), &prm); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("failed to unmarshal permissions"))
	}

	e, err := h.auth.CreateEntity(ctx, authSec, entitySec, prm, req.Msg.Note)
	if errors.Is(err, authservice.ErrUnauthorized) {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	} else if errors.Is(err, authservice.ErrInvalidArg{}) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	} else if err != nil {
		h.l.Error().Err(err).Msg("failed to create an auth entity")
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	return connect.NewResponse(&v1.CreateEntityResponse{Id: e.Id}), nil
}

func (h *Handler) Login(
	ctx context.Context,
	req *connect.Request[v1.LoginRequest],
) (*connect.Response[v1.LoginResponse], error) {
	tok, err := h.auth.Login(ctx, req.Msg.Id, req.Msg.Secret)

	if errors.Is(err, authservice.ErrUnauthorized) {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	} else if err != nil {
		h.l.Error().Err(err).Msg("failed to perform login operation")
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	return connect.NewResponse(&v1.LoginResponse{Token: tok}), nil
}

func (h *Handler) Logout(
	context.Context,
	*connect.Request[v1.LogoutRequest],
) (*connect.Response[v1.LogoutResponse], error) {
	return nil, nil
}
