package config

import "time"

type RedisConfig struct {
	Addr        string        `yaml:"addr" env:"REDIS_ADDR"`
	Password    string        `yaml:"password" env:"REDIS_PASSWORD"`
	User        string        `yaml:"user" env:"REDIS_USER"`
	DB          int           `yaml:"db" env:"REDIS_DB"`
	MaxRetries  int           `yaml:"max_retries" env:"REDIS_MAX_RETRIES"`
	DialTimeout time.Duration `yaml:"dial_timeout" env:"REDIS_DIAL_TIMEOUT"`
	Timeout     time.Duration `yaml:"timeout" env:"REDIS_TIMEOUT"`
}