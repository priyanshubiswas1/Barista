package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port  string `json:"PORT" default:"8080"`
	DbUrl string `json:"DB_URL" default:"postgres://admin:7866@shop_db:5432/shop?sslmode=disable"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
