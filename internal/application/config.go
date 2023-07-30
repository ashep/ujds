package application

import (
	_ "embed"

	"github.com/ashep/ujds/internal/server"
)

//go:embed schema.json
var Schema []byte

type Database struct {
	DSN string `json:"dsn,omitempty" yaml:"dsn,omitempty"`
}

type Config struct {
	DB     Database      `json:"db,omitempty" yaml:"db,omitempty"`
	Server server.Config `json:"server,omitempty" yaml:"server,omitempty"`
}
