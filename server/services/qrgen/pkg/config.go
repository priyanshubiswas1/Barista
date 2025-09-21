package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RMQUsername       string `json:"RMQ_USER_NAME" default:"admin"`
	RMQPassword       string `json:"RMQ_PASSWORD" default:"7866"`
	RMQHost           string `json:"RMQ_HOST" default:"rabbitmqservice:5672"`
	RMQVhost          string `json:"RMQ_VHOST" default:"barista"`
	ConsumeQueueName  string `json:"QUEUE_NAME" default:"qr_gen_queue"`
	ConsumeBindingKey string `json:"CON_BINDING_KEY" default:"gen_qr"`
	Exchange          string `json:"EXCHANGE" default:"gen_qr_event"`
	ConsumerName      string `json:"CONSUMER_NAME" default:"qr_gen_service"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
