package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port        string `json:"PORT" default:"8080"`
	DbUrl       string `json:"DB_URL" default:"postgres://admin:7866@session_db:5432/session?sslmode=disable"`
	RedisUrl    string `json:"REDIS_URL" default:"redis:6379"`
	TableClient string `envconfig:"TABLE_CLIENT" default:"localhost:8083"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
