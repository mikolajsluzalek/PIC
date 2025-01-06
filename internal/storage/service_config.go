package storage

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

func readConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse config")
	}

	return cfg, nil
}
