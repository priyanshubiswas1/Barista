package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port         string `envconfig:"PORT"`
	DbUrl        string `envconfig:"DB_URL"`
	RMQUsername  string `envconfig:"RMQ_USER_NAME"`
	RMQPassword  string `envconfig:"RMQ_PASSWORD"`
	RMQHost      string `envconfig:"RMQ_HOST"`
	RMQVhost     string `envconfig:"RMQ_VHOST"`
	ClientUrl    string `envconfig:"CLIENT_URL"`
	Exchange     string `envconfig:"EXCHANGE"`
	ConsumerName string `envconfig:"CONSUMER_NAME"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
