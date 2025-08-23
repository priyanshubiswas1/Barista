package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port       string `envconfig:"PORT" default:"8000"`
	AuthClient string `envconfig:"AUTH_CLIENT" default:"localhost:8080"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
