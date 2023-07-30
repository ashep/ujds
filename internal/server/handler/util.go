package handler

import (
	"errors"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
)

func grpcErr(err error, proc, msg string, l zerolog.Logger) error {
	switch {
	case errors.As(err, &apperrors.NotFoundError{}):
		l.Info().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeNotFound, err)
	case errors.As(err, &apperrors.InvalidArgError{}), errors.As(err, &apperrors.EmptyArgError{}):
		l.Info().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.As(err, &apperrors.AlreadyExistsError{}):
		l.Info().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeAlreadyExists, err)
	case errors.As(err, &apperrors.AccessDeniedError{}):
		l.Info().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeUnauthenticated, err)
	default:
		l.Error().Str("proc", proc).Err(err).Msg(msg)
		return connect.NewError(connect.CodeInternal, nil)
	}
}
