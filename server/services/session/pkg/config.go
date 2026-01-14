package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port           string `envconfig:"PORT"`
	DbUrl          string `envconfig:"DB_URL"`
	RedisUrl       string `envconfig:"REDIS_URL"`
	TableService   string `envconfig:"TABLE_SERVICE_URL"`
	HistoryService string `envconfig:"HISTORY_SERVICE_URL"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
