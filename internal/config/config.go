package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Brokers []string `yaml:"brokers" env-required:"true"`
	Topics  []string `yaml:"topics" env-required:"true"`
	Group   string   `yaml:"consumer_group" env-required:"true"`
}

func New() (*Config, error) {
	const op = "config.New"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("%s: failed to get CONFIG_PATH env var", op)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cfg, nil
}
