package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port   string `env:"PORT" env-default:":8080"`
	DBPath string `env:"DB_PATH" env-required:"true"`
}

func GetConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatal("Ошибка чтения конфигурации")
	}
	
	return cfg
}