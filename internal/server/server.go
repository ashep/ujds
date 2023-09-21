package server

import (
	"context"
	"errors"
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
	c Config
	s *http.Server
	l zerolog.Logger
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

	interceptors := connect.WithInterceptors(
		NewAuthInterceptor(cfg.AuthToken),
	)

	mux := http.NewServeMux()
	mux.Handle(indexconnect.NewIndexServiceHandler(ih, interceptors))
	mux.Handle(recordconnect.NewRecordServiceHandler(rh, interceptors))

	return &Server{
		c: cfg,
		s: &http.Server{Addr: cfg.Address, Handler: cors(mux), ReadTimeout: readTimeout},
		l: l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		if errF := s.s.Close(); errF != nil {
			s.l.Error().Err(errF).Msg("server close failed")
		}
	}()

	s.l.Info().Str("addr", s.c.Address).Msg("server is starting")

	if s.c.AuthToken == "" {
		s.l.Warn().Msg("empty auth token, should be used only for development purposes")
	}

	if err := s.s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve failed: %w", err)
	}

	s.l.Info().Str("addr", s.c.Address).Msg("server stopped")

	return nil
}
