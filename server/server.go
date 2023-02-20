package server

import (
	"context"
	"net/http"

	"github.com/ashep/datapimp/dataservice"
	"github.com/ashep/datapimp/gen/proto/datapimp/v1/datapimpv1connect"
	"github.com/rs/zerolog"
)

const (
	exchangeName = "amq.fanout"
	queueName    = "items"
)

type Server struct {
	cfg Config
	ds  *dataservice.Service
	l   zerolog.Logger
}

func New(cfg Config, ds *dataservice.Service, l zerolog.Logger) *Server {
	if cfg.Addr == "" {
		cfg.Addr = "localhost:8080"
	}

	return &Server{
		cfg: cfg,
		ds:  ds,
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	p, h := datapimpv1connect.NewItemServiceHandler(s)
	mux := http.NewServeMux()
	mux.Handle(p, h)
	srv := &http.Server{Addr: s.cfg.Addr, Handler: mux}

	go func() {
		<-ctx.Done()
		if errF := srv.Close(); errF != nil {
			s.l.Error().Err(errF).Msg("failed to close server")
		}
	}()

	s.l.Debug().Str("addr", s.cfg.Addr).Msg("starting server")
	return srv.ListenAndServe()
}
