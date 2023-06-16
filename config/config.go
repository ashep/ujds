package config

import (
	_ "embed"

	"github.com/ashep/ujds/api"
	"github.com/ashep/ujds/server"
)

//go:embed schema.json
var Schema []byte

type Database struct {
	DSN string `json:"dsn,omitempty" yaml:"dsn,omitempty"`
}

type Config struct {
	DB     Database      `json:"db,omitempty" yaml:"db,omitempty"`
	API    api.Config    `json:"api,omitempty" yaml:"api,omitempty"`
	Server server.Config `json:"server,omitempty" yaml:"server,omitempty"`
}
