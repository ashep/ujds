package handler

import (
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/api"
)

type Handler struct {
	api *api.API
	l   zerolog.Logger
}

func New(api *api.API, l zerolog.Logger) *Handler {
	return &Handler{
		api: api,
		l:   l,
	}
}
