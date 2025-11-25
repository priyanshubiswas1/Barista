package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port          string `envconfig:"PORT" default:"8080"`
	DbUrl         string `envconfig:"DB_URL" default:"postgresql://postgres.qwjlqyuxcadqwvmarwlo:Barista@2025@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"`
	RedisUrl      string `envconfig:"REDIS_URL" default:"redis:6379"`
	TableClient   string `envconfig:"TABLE_CLIENT" default:"localhost:8083"`
	HistoryClient string `envconfig:"HISTORY_CLIENT" default:"localhost:8087"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
