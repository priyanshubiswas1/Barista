package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port             string `envconfig:"PORT"`
	DbUrl            string `envconfig:"DB_URL"`
	SessionService   string `envconfig:"SESSION_SERVICE_URL"`
	WebSocketService string `envconfig:"WEBSOCKET_SERVICE_URL"`
	HistoryService   string `envconfig:"HISTORY_SERVICE_URL"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
