package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	API_KEY string `env:"API_KEY"`
	Port    int    `env:"PORT"`
}

func LoadConfig() *Config {
	path := "D:/currency-exchange/exchange_service/config/.env"
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("Cannot read config: " + err.Error())
	}

	return &config
}
