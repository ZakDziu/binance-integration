package config

import (
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

type Configs struct {
	Server ServerConfig
}

type ServerConfig struct {
	ServerPort  string   `env:"SERVER_PORT"`
	ReadTimeout Duration `env:"READ_TIMEOUT"`
}

func New() (*Configs, error) {
	var config Configs
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
