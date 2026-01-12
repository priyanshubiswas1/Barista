package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                string `envconfig:"PORT"`
	DbUrl               string `envconfig:"DB_URL"`
	CloudinaryUrl       string `envconfig:"CLOUDINARY_URL"`
	CloudinaryApiKey    string `envconfig:"CLOUDINARY_API_KEY"`
	CloudinarySecret    string `envconfig:"CLOUDINARY_API_SECRET"`
	CloudinaryCloudName string `envconfig:"CLOUDINARY_CLOUD_NAME"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
