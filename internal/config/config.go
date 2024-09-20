package config

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

func Load[C any]() (C, error) {
	var cfg C
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		return cfg, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}
