package client

import (
	"context"

	"github.com/bufbuild/connect-go"
)

func NewAuthInterceptor(apiKey string) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			req.Header().Set("Authorization", "Bearer "+apiKey)
			return next(ctx, req)
		}
	}

	return interceptor
}
