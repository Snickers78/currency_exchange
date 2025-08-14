package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `env:"env" env-default:"local"`
	StoragePath string        `env:"storage_path" env-required:"true"`
	TockenTTL   time.Duration `env:"tocken_ttl" env-required:"true"`
	Secret      string        `env:"secret" env-required:"true"`
	Port        int           `env:"port"`
	Timeout     time.Duration `env:"timeout"`
	Broker1     string        `env:"broker1"`
	Broker2     string        `env:"broker2"`
}

func MustLoad() *Config {
	path := "./config/.env"
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("Cannot read config: " + err.Error())
	}

	return &config
}
