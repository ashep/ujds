package handler

import (
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/apperrors"
)

func grpcErr(err error, proc, msg string, l zerolog.Logger) error {
	switch {
	case errors.Is(err, apperrors.NotFoundError{}):
		l.Info().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, apperrors.InvalidArgError{}), errors.Is(err, apperrors.EmptyArgError{}):
		l.Info().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, apperrors.AlreadyExistsError{}):
		l.Info().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeAlreadyExists, err)
	case errors.Is(err, apperrors.AccessDeniedError{}):
		l.Info().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeUnauthenticated, err)
	default:
		l.Error().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeInternal, nil)
	}
}
