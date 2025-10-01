package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port       string `json:"PORT" default:"8080"`
	DbUrl      string `json:"DB_URL" default:"postgres://admin:7866@user_db:5432/user?sslmode=disable"`
	ShopClient string `envconfig:"SHOP_CLIENT" default:"localhost:8081"`
	MenuClient string `envconfig:"MENU_CLIENT" default:"localhost:8082"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
