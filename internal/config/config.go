package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port           string      `env:"PORT" env-default:":8080"`
	DBPath         string      `env:"DB_PATH" env-required:"true"`
	OperatorAPIKey string      `env:"OPERATOR_API_KEY" env-required:"true"`
	Redis          RedisConfig `yaml:"redis"`
	WebhookURL     string      `env:"WEBHOOK_URL" env-required:"true"`
	WindowMin      int         `env:"STATS_TIME_WINDOW_MINUTES" env-required:"true"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	return &cfg, nil
}
