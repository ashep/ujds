package server

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/sdk/proto/ujds/v1/v1connect"
	"github.com/ashep/ujds/server/handler"
)

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
	mux := http.NewServeMux()

	hdl := handler.New(s.api, s.l)
	mux.Handle(v1connect.NewSchemaServiceHandler(hdl))
	mux.Handle(v1connect.NewRecordServiceHandler(hdl))

	srv := &http.Server{Addr: s.cfg.Address, Handler: mux}

	go func() {
		<-ctx.Done()
		if errF := srv.Close(); errF != nil {
			s.l.Error().Err(errF).Msg("failed to close server")
		}
	}()

	s.l.Info().Str("addr", s.cfg.Address).Msg("starting server")
	return srv.ListenAndServe()
}
