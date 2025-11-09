package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          string `envconfig:"PORT" default:"8000"`
	AuthClient    string `envconfig:"AUTH_CLIENT" default:"localhost:8080"`
	ShopClient    string `envconfig:"SHOP_CLIENT" default:"localhost:8081"`
	MenuClient    string `envconfig:"MENU_CLIENT" default:"localhost:8082"`
	TableClient   string `envconfig:"TABLE_CLIENT" default:"localhost:8083"`
	SessionClient string `envconfig:"SESSION_CLIENT" default:"localhost:8084"`
	OrderClient   string `envconfig:"ORDER_CLIENT" default:"localhost:8085"`
	PaymentClient string `envconfig:"PAYMENT_CLIENT" default:"localhost:8086"`
	HistoryClient string `envconfig:"HISTORY_CLIENT" default:"localhost:8087"`
	JwtKey        string `envconfig:"JWT_KEY" default:"key"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
