package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"

	indexconnect "github.com/ashep/ujds/sdk/proto/ujds/index/v1/v1connect"
	recordconnect "github.com/ashep/ujds/sdk/proto/ujds/record/v1/v1connect"
)

const readTimeout = time.Second * 5

type Server struct {
	cfg Config
	srv *http.Server
	l   zerolog.Logger
}

func New(
	cfg Config,
	ih indexconnect.IndexServiceHandler,
	rh recordconnect.RecordServiceHandler,
	l zerolog.Logger,
) *Server {
	if cfg.Address == "" {
		cfg.Address = ":9000"
	}

	interceptors := connect.WithInterceptors(NewAuthInterceptor(cfg.AuthToken))

	mux := http.NewServeMux()
	mux.Handle(indexconnect.NewIndexServiceHandler(ih, interceptors))
	mux.Handle(recordconnect.NewRecordServiceHandler(rh, interceptors))

	srv := &http.Server{Addr: cfg.Address, Handler: mux, ReadTimeout: readTimeout}

	return &Server{cfg: cfg, srv: srv, l: l}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		if errF := s.srv.Close(); errF != nil {
			s.l.Error().Err(errF).Msg("failed to close server")
		}
	}()

	s.l.Info().Str("addr", s.cfg.Address).Msg("starting server")

	if err := s.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("ListenAndServe failed: %w", err)
	}

	return nil
}
