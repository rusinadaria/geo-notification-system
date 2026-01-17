package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Port           string      `env:"PORT" env-default:":8080"`
	DBPath         string      `env:"DB_PATH" env-required:"true"`
	OperatorAPIKey string      `env:"OPERATOR_API_KEY" env-required:"true"`
	Redis          RedisConfig `yaml:"redis"`
	WebhookURL     string      `env:"WEBHOOK_URL" env-required:"true"`
}

func GetConfig() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal("Ошибка чтения конфигурации")
	}

	return &cfg
}
