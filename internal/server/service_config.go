package server

import (
	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JWTSecret string `envconfig:"JWT_SECRET" required:"true"`
}

func readConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, errors.Errorf("failed to parse config; error=%v", err)
	}

	return cfg, nil
}
