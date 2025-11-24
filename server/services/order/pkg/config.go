package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port             string `json:"PORT" default:"8080"`
	DbUrl            string `json:"DB_URL" default:"postgresql://postgres.qwjlqyuxcadqwvmarwlo:Barista@2025@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"`
	SessionClient    string `envconfig:"SESSION_CLIENT" default:"localhost:8084"`
	WebSocketService string `envconfig:"WEBSOCKET_URL" default:"localhost:8087"`
	HistoryClient    string `envconfig:"HISTORY_CLIENT" default:"localhost:8087"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
