package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port           string `envconfig:"PORT"`
	AuthService    string `envconfig:"AUTH_CLIENT"`
	ShopService    string `envconfig:"SHOP_CLIENT"`
	MenuService    string `envconfig:"MENU_CLIENT"`
	TableService   string `envconfig:"TABLE_CLIENT"`
	SessionService string `envconfig:"SESSION_CLIENT"`
	OrderService   string `envconfig:"ORDER_CLIENT"`
	PaymentService string `envconfig:"PAYMENT_CLIENT"`
	HistoryService string `envconfig:"HISTORY_CLIENT"`
	JwtKey         string `envconfig:"JWT_KEY"`
	AllowedOrigins string `envconfig:"ALLOWED_ORIGINS"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
