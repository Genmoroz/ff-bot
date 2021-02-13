package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken string `split_words:"true" required:"true"`
}

func ReadEnv() (Config, error) {
	var config Config
	if err := envconfig.Process("APP", &config); err != nil {
		return Config{}, fmt.Errorf("failed to read envs: %w", err)
	}

	return config, nil
}
