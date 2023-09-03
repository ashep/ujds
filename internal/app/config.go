package app

import (
	"github.com/ashep/ujds/internal/server"
)

type Database struct {
	DSN string `json:"dsn" yaml:"dsn"`
}

type Config struct {
	DB     Database      `json:"db" yaml:"db"`
	Server server.Config `json:"server" yaml:"server"`
}
