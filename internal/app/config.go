package app

type Server struct {
	Addr      string `json:"addr" yaml:"addr"`
	AuthToken string `json:"auth_token" yaml:"auth_token"`
}

type Database struct {
	DSN string `json:"dsn" yaml:"dsn"`
}

type Validation struct {
	Record string `json:"record" yaml:"record"`
}

type Config struct {
	DB         Database   `json:"db" yaml:"db"`
	Server     Server     `json:"server" yaml:"server"`
	Validation Validation `json:"validation" yaml:"validation"`
}
