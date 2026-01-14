package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port           string `envconfig:"PORT"`
	AuthService    string `envconfig:"AUTH_SERVICE_URL"`
	ShopService    string `envconfig:"SHOP_SERVICE_URL"`
	MenuService    string `envconfig:"MENU_SERVICE_URL"`
	TableService   string `envconfig:"TABLE_SERVICE_URL"`
	SessionService string `envconfig:"SESSION_SERVICE_URL"`
	OrderService   string `envconfig:"ORDER_SERVICE_URL"`
	PaymentService string `envconfig:"PAYMENT_SERVICE_URL"`
	HistoryService string `envconfig:"HISTORY_SERVICE_URL"`
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
