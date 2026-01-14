package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port              string `envconfig:"PORT"`
	RMQUsername       string `envconfig:"RMQ_USER_NAME"`
	RMQPassword       string `envconfig:"RMQ_PASSWORD"`
	RMQHost           string `envconfig:"RMQ_HOST"`
	RMQVhost          string `envconfig:"RMQ_VHOST"`
	ConsumeQueueName  string `envconfig:"QUEUE_NAME"`
	ConsumeBindingKey string `envconfig:"CON_BINDING_KEY"`
	Exchange          string `envconfig:"EXCHANGE"`
	ConsumerName      string `envconfig:"CONSUMER_NAME"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
