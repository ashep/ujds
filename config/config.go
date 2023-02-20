package config

import (
	"io"
	"os"

	"github.com/ashep/datapimp/dataservice"
	"github.com/ashep/datapimp/mq"
	"github.com/ashep/datapimp/server"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
	MQ     mq.Config          `yaml:"mq"`
	Data   dataservice.Config `yaml:"data"`
	Server server.Config      `yaml:"server"`
}

func Parse(in []byte) (Config, error) {
	r := Config{}
	err := yaml.Unmarshal(in, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func ParseFromPath(path string) (Config, error) {
	fp, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer fp.Close()

	b, err := io.ReadAll(fp)
	if err != nil {
		return Config{}, err
	}

	return Parse(b)
}
