package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type EnvConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

// GetEnvConfig parses and returns environment variables
func GetEnvConfig() EnvConfig {
	var cfg EnvConfig

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current environment variables:", cfg)
	return cfg
}
