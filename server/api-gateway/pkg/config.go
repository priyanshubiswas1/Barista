package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port           string `envconfig:"PORT"`
	AuthClient     string `envconfig:"AUTH_CLIENT"`
	ShopClient     string `envconfig:"SHOP_CLIENT"`
	MenuClient     string `envconfig:"MENU_CLIENT"`
	TableClient    string `envconfig:"TABLE_CLIENT"`
	SessionClient  string `envconfig:"SESSION_CLIENT"`
	OrderClient    string `envconfig:"ORDER_CLIENT"`
	PaymentClient  string `envconfig:"PAYMENT_CLIENT"`
	HistoryClient  string `envconfig:"HISTORY_CLIENT"`
	JwtKey         string `envconfig:"JWT_KEY"`
	AllowedOrigins string `envconfig:"ALLOWED_ORIGINS"`

	// Garage (S3) for 3D scan downloads
	GarageEndpoint  string `envconfig:"GARAGE_ENDPOINT"`
	GarageRegion    string `envconfig:"GARAGE_REGION"`
	GarageAccessKey string `envconfig:"GARAGE_ACCESS_KEY"`
	GarageSecretKey string `envconfig:"GARAGE_SECRET_KEY"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
