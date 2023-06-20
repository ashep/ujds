package server

type Config struct {
	Address   string `json:"address,omitempty" yaml:"address,omitempty"`
	AuthToken string `json:"auth_token,omitempty" yaml:"auth_token,omitempty"`
}
