package config

import (
	env "github.com/d4rthvadr/dusky-go/internal/env"
)

type ServerConfig struct {
	Port string
	Host string
}

type AppConfig struct {
	Server ServerConfig
}

func InitializeConfig() (*AppConfig, error) {

	serverAddr := env.GetEnv("ADDR", "3000")

	config := &AppConfig{
		Server: ServerConfig{
			Host: serverAddr,
		},
	}
	return config, nil
}
