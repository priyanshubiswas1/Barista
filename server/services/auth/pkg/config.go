package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        string `envconfig:"PORT"`
	DbUrl       string `envconfig:"DB_URL"`
	ShopService string `envconfig:"SHOP_CLIENT"`
	MenuService string `envconfig:"MENU_CLIENT"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
