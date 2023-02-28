package interceptor

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
)

func Auth(l zerolog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			basicStr := strings.TrimPrefix(req.Header().Get("Authorization"), "Basic ")
			if basicStr != "" {
				authB, err := base64.StdEncoding.DecodeString(basicStr)
				if err != nil {
					l.Warn().Err(err).Msg("failed to decode authorization header")
					return nil, connect.NewError(connect.CodeUnauthenticated, nil)
				}
				basicStr = string(authB)

				basicSplit := strings.Split(basicStr, ":")
				if len(basicSplit) != 2 {
					return nil, connect.NewError(connect.CodeUnauthenticated, nil)
				}

				ctx = context.WithValue(ctx, "authId", basicSplit[0])
				ctx = context.WithValue(ctx, "authSecret", basicSplit[1])
			}

			bearerStr := strings.TrimPrefix(req.Header().Get("Authorization"), "Bearer ")
			if bearerStr != "" {
				ctx = context.WithValue(ctx, "authToken", bearerStr)
			}

			return next(ctx, req)
		}
	}
}
