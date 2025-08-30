package app

type Server struct {
	AuthToken string `json:"auth_token" yaml:"auth_token"`
}

type Database struct {
	DSN string `json:"dsn" yaml:"dsn"`
}

type Config struct {
	DB     Database `json:"db" yaml:"db"`
	Server Server   `json:"server" yaml:"server"`
}
