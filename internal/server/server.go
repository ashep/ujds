package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/api"
	"github.com/ashep/ujds/internal/server/handler"
	"github.com/ashep/ujds/sdk/proto/ujds/v1/v1connect"
)

const readTimeout = time.Second * 5

type Server struct {
	cfg Config
	api *api.API
	l   zerolog.Logger
}

func New(cfg Config, api *api.API, l zerolog.Logger) *Server {
	if cfg.Address == "" {
		cfg.Address = ":9000"
	}

	return &Server{
		cfg: cfg,
		api: api,
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	interceptors := connect.WithInterceptors(NewAuthInterceptor(s.cfg.AuthToken))
	mux := http.NewServeMux()

	hdl := handler.New(s.api, s.l)
	mux.Handle(v1connect.NewIndexServiceHandler(hdl, interceptors))
	mux.Handle(v1connect.NewRecordServiceHandler(hdl, interceptors))

	srv := &http.Server{
		Addr:        s.cfg.Address,
		Handler:     mux,
		ReadTimeout: readTimeout,
	}

	go func() {
		<-ctx.Done()

		if errF := srv.Close(); errF != nil {
			s.l.Error().Err(errF).Msg("failed to close server")
		}
	}()

	s.l.Info().Str("addr", s.cfg.Address).Msg("starting server")

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("ListenAndServe failed: %w", err)
	}

	return nil
}
