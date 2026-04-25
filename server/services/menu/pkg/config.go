package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                string `envconfig:"PORT"`
	DbUrl               string `envconfig:"DB_URL"`
	CloudinaryUrl       string `envconfig:"CLOUDINARY_URL"`
	CloudinaryApiKey    string `envconfig:"CLOUDINARY_API_KEY"`
	CloudinarySecret    string `envconfig:"CLOUDINARY_API_SECRET"`
	CloudinaryCloudName string `envconfig:"CLOUDINARY_CLOUD_NAME"`

	// RabbitMQ (3D scan jobs)
	RMQUsername    string `envconfig:"RMQ_USER_NAME"`
	RMQPassword    string `envconfig:"RMQ_PASSWORD"`
	RMQHost        string `envconfig:"RMQ_HOST"`
	RMQVhost       string `envconfig:"RMQ_VHOST"`
	ScanExchange   string `envconfig:"SCAN_EXCHANGE"`
	ScanQueueName  string `envconfig:"SCAN_QUEUE_NAME"`
	ScanRoutingKey string `envconfig:"SCAN_ROUTING_KEY"`
	ConsumerName   string `envconfig:"CONSUMER_NAME"`

	// Garage (S3)
	GarageEndpoint        string `envconfig:"GARAGE_ENDPOINT"`
	GarageRegion          string `envconfig:"GARAGE_REGION"`
	GarageAccessKey       string `envconfig:"GARAGE_ACCESS_KEY"`
	GarageSecretKey       string `envconfig:"GARAGE_SECRET_KEY"`
	GarageBucket          string `envconfig:"GARAGE_BUCKET"`
	GaragePublicURLPrefix string `envconfig:"GARAGE_PUBLIC_URL_PREFIX"`

	// Hunyuan3D runner
	HunyuanPython string `envconfig:"HUNYUAN_PYTHON"`
	HunyuanScript string `envconfig:"HUNYUAN_SCRIPT"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
