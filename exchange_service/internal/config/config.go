package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	API_KEY string        `env:"API_KEY"`
	Port    int           `env:"PORT"`
	Timeout time.Duration `env:"TIMEOUT"`
	Broker1 string        `env:"BROKER1"`
	Broker2 string        `env:"BROKER2"`
}

func LoadConfig() *Config {
	// configPath := flag.String("config", "", "Path to config file")
	// flag.Parse()

	path := "./config/.env"
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("Cannot read config: " + err.Error())
	}

	return &config
}
