package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AuthServicePort     int           `env:"AUTH_PORT"`
	ExchangeServicePort int           `env:"EXCHANGE_PORT"`
	GatewayPort         int           `env:"GATEWAY_PORT"`
	Timeout             time.Duration `env:"TIMEOUT"`
	Secret              string        `env:"SECRET" env-required:"true"`
}

func LoadConfig() *Config {
	path := "./config/.env"
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("Cannot read config: " + err.Error())
	}

	return &config
}
