package pkg

import (
	"os"
)

type Config struct {
	Port       string
	AuthClient string
}

func GetConfig() (*Config, error) {
	return &Config{
		Port:       os.Getenv("PORT"),
		AuthClient: os.Getenv("AUTH_CLIENT"),
	}, nil
}
