package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port         string `json:"PORT" default:"8080"`
	DbUrl        string `json:"DB_URL" default:"postgres://admin:7866@table_db:5432/table?sslmode=disable"`
	RMQUsername  string `json:"RMQ_USER_NAME" default:"admin"`
	RMQPassword  string `json:"RMQ_PASSWORD" default:"7866"`
	RMQHost      string `json:"RMQ_HOST" default:"rabbitmqservice:5672"`
	RMQVhost     string `json:"RMQ_VHOST" default:"barista"`
	ClientUrl    string `json:"CLIENT_URL" default:"http://clienturl"`
	Exchange     string `json:"EXCHANGE" default:"gen_qr_event"`
	ConsumerName string `json:"CONSUMER_NAME" default:"qr_gen_service"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
