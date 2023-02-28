package server

import (
	"context"
	"net/http"

	"github.com/ashep/datapimp/authservice"
	"github.com/ashep/datapimp/dataservice"
	"github.com/ashep/datapimp/gen/proto/datapimp/v1/v1connect"
	"github.com/ashep/datapimp/server/handler"
	"github.com/ashep/datapimp/server/interceptor"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
)

const (
	exchangeName = "amq.fanout"
	queueName    = "items"
)

type Server struct {
	cfg  Config
	auth *authservice.Service
	data *dataservice.Service
	l    zerolog.Logger
}

func New(cfg Config, auth *authservice.Service, data *dataservice.Service, l zerolog.Logger) *Server {
	if cfg.Addr == "" {
		cfg.Addr = "localhost:8080"
	}

	return &Server{
		cfg:  cfg,
		auth: auth,
		data: data,
		l:    l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()

	hdl := handler.New(s.auth, s.data, s.l)

	interceptors := connect.WithInterceptors(interceptor.Auth(s.l))

	p, h := v1connect.NewAuthServiceHandler(hdl, interceptors)
	mux.Handle(p, h)

	p, h = v1connect.NewDataServiceHandler(hdl, interceptors)
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
