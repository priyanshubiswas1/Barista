package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port         string `envconfig:"PORT" default:"8080"`
	DbUrl        string `envconfig:"DB_URL" default:"postgresql://postgres.qwjlqyuxcadqwvmarwlo:Barista@2025@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"`
	RMQUsername  string `envconfig:"RMQ_USER_NAME" default:"admin"`
	RMQPassword  string `envconfig:"RMQ_PASSWORD" default:"7866"`
	RMQHost      string `envconfig:"RMQ_HOST" default:"rabbitmqservice:5672"`
	RMQVhost     string `envconfig:"RMQ_VHOST" default:"barista"`
	ClientUrl    string `envconfig:"CLIENT_URL" default:"https://assembly-blackberry-correct-penny.trycloudflare.com"`
	Exchange     string `envconfig:"EXCHANGE" default:"gen_qr_event"`
	ConsumerName string `envconfig:"CONSUMER_NAME" default:"qr_gen_service"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
