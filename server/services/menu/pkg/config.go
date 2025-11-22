package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                string `json:"PORT" default:"8080"`
	DbUrl               string `json:"DB_URL" default:"postgres://admin:7866@menu_db:5432/menu?sslmode=disable"`
	CloudinaryUrl       string `json:"CLOUDINARY_URL" default:"cloudinary://759689914359316:HrHhR6oR09cGl2Ele49DeAql8k8@dw71avwb7"`
	CloudinaryApiKey    string `json:"CLOUDINARY_API_KEY" default:"759689914359316"`
	CloudinarySecret    string `json:"CLOUDINARY_API_SECRET" default:"HrHhR6oR09cGl2Ele49DeAql8k8"`
	CloudinaryCloudName string `json:"CLOUDINARY_CLOUD_NAME" default:"dw71avwb7"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
