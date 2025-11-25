package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RMQUsername       string `envconfig:"RMQ_USER_NAME" default:"admin"`
	RMQPassword       string `envconfig:"RMQ_PASSWORD" default:"7866"`
	RMQHost           string `envconfig:"RMQ_HOST" default:"rabbitmqservice:5672"`
	RMQVhost          string `envconfig:"RMQ_VHOST" default:"barista"`
	ConsumeQueueName  string `envconfig:"QUEUE_NAME" default:"qr_gen_queue"`
	ConsumeBindingKey string `envconfig:"CON_BINDING_KEY" default:"gen_qr"`
	Exchange          string `envconfig:"EXCHANGE" default:"gen_qr_event"`
	ConsumerName      string `envconfig:"CONSUMER_NAME" default:"qr_gen_service"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
