package server

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
)

func NewAuthInterceptor(token string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if token == "" {
				return next(ctx, req)
			}

			if token == strings.ReplaceAll(req.Header().Get("Authorization"), "Bearer ", "") {
				return next(ctx, req)
			}

			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not authorized"))
		}
	}
}
