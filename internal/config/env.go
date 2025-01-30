package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
)

type EnvConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

// GetEnvConfig parses and returns environment variables
func GetEnvConfig() EnvConfig {
	var cfg EnvConfig

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Current environment variables: %+v", cfg)
	return cfg
}
