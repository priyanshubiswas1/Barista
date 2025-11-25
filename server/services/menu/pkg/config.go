package pkg

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                string `envconfig:"PORT" default:"8080"`
	DbUrl               string `envconfig:"DB_URL" default:"postgresql://postgres.qwjlqyuxcadqwvmarwlo:Barista@2025@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"`
	CloudinaryUrl       string `envconfig:"CLOUDINARY_URL" default:"cloudinary://759689914359316:HrHhR6oR09cGl2Ele49DeAql8k8@dw71avwb7"`
	CloudinaryApiKey    string `envconfig:"CLOUDINARY_API_KEY" default:"759689914359316"`
	CloudinarySecret    string `envconfig:"CLOUDINARY_API_SECRET" default:"HrHhR6oR09cGl2Ele49DeAql8k8"`
	CloudinaryCloudName string `envconfig:"CLOUDINARY_CLOUD_NAME" default:"dw71avwb7"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
