package server

import "efaturas-xtreme/pkg/validator"

type Config struct {
	API struct {
		Name         string `env:"API_NAME,default=efaturas-extreme" validate:"notblank"`
		Port         uint   `env:"API_PORT,default=8080" validate:"notblank,number"`
		Version      string `env:"API_VERSION,default=v1" validate:"notblank"`
		Environment  string `env:"API_ENVIRONMENT,default=LOCAL" validate:"notblank,oneof=LOCAL DEV PROD"`
		Certificates string `env:"API_CERTIFICATES,default=certs/" validate:"notblank"`
	}

	Database struct {
		URI  string `env:"DATABASE_URI,default=mongodb://localhost:27017"`
		Name string `env:"DATABASE_NAME,default=efaturas"`
	}
}

func (cfg *Config) Validate() error {
	return validator.Validate(cfg)
}
