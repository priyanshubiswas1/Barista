package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port  string `json:"PORT" default:"8080"`
	DbUrl string `json:"DB_URL" default:"postgresql://postgres.qwjlqyuxcadqwvmarwlo:Barista@2025@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
