package server

type Config struct {
	Address   string `json:"address" yaml:"address"`
	AuthToken string `json:"auth_token" yaml:"auth_token"`
}
