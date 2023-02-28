package handler

import (
	"github.com/ashep/datapimp/authservice"
	"github.com/ashep/datapimp/dataservice"
	"github.com/rs/zerolog"
)

type Handler struct {
	auth *authservice.Service
	data *dataservice.Service
	l    zerolog.Logger
}

func New(auth *authservice.Service, data *dataservice.Service, l zerolog.Logger) *Handler {
	return &Handler{
		auth: auth,
		data: data,
		l:    l,
	}
}
