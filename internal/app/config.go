package app

import (
	"github.com/ashep/ujds/internal/server"
)

type LogServer struct {
	URL      string `json:"url" yaml:"url"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type MetricServer struct {
	Addr string `json:"addr" yaml:"addr" default:":9090"`
	Path string `json:"path" yaml:"path" default:"/metrics"`
}

type Database struct {
	DSN string `json:"dsn" yaml:"dsn"`
}

type Config struct {
	LogServer    LogServer     `json:"log_server" yaml:"log_server"`
	MetricServer MetricServer  `json:"metric_server" yaml:"metric_server"`
	DB           Database      `json:"db" yaml:"db"`
	Server       server.Config `json:"server" yaml:"server"`
}
