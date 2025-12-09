package app

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	Addr      string `json:"addr" yaml:"addr"`
	AuthToken string `json:"auth_token" yaml:"auth_token"`
}

type Database struct {
	DSN string `json:"dsn" yaml:"dsn"`
}

type Validation struct {
	Index       string                     // to load from env var
	IndexStruct map[string]json.RawMessage `json:"index" yaml:"index" env:"ignore"`
}

type Config struct {
	DB         Database   `json:"db" yaml:"db"`
	Server     Server     `json:"server" yaml:"server"`
	Validation Validation `json:"validation" yaml:"validation"`
}

func (c *Config) Validate() error {
	if c.Validation.Index != "" {
		if err := json.Unmarshal([]byte(c.Validation.Index), &c.Validation.IndexStruct); err != nil {
			return fmt.Errorf("VALIDATION_INDEX: parse JSON: %w", err)
		}
	}

	if c.Validation.IndexStruct == nil {
		c.Validation.IndexStruct = make(map[string]json.RawMessage)
	}

	c.Validation.IndexStruct[".*"] = json.RawMessage(`{}`) // validate all for valid JSON

	return nil
}
