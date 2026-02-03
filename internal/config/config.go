package config

import (
	"time"

	env "github.com/d4rthvadr/dusky-go/internal/env"
)

type ServerConfig struct {
	Port string
	Host string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type AppConfig struct {
	Server ServerConfig
	Db     DbConfig
}

func InitializeConfig() (*AppConfig, error) {

	serverAddr := env.GetEnv("ADDR", "3000")
	dbAddr := env.GetEnv("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/duskydb?sslmode=disable")
	maxOpenConns := env.GetEnvAsInt("DB_MAX_OPEN_CONNS", 30)
	maxIdleConns := env.GetEnvAsInt("DB_MAX_IDLE_CONNS", 30)
	maxIdleTime := env.GetEnvAsDuration("DB_MAX_IDLE_TIME", time.Minute*15)

	config := &AppConfig{
		Server: ServerConfig{
			Host: serverAddr,
		},
		Db: DbConfig{
			Addr:         dbAddr,
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			MaxIdleTime:  maxIdleTime,
		},
	}
	return config, nil
}
