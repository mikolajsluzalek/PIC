package api

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	JWTSecret string `envconfig:"JWT_SECRET" required:"true"`
}

func readConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse config")
	}

	return cfg, nil
}
