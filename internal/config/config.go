package config

import (
	"time"

	env "github.com/d4rthvadr/dusky-go/internal/utils"
)

type serverConfig struct {
	Port string
	Host string
}

type dbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type sendGridConfig struct {
	APIKey string
}

type MailConfig struct {
	Expiry    time.Duration
	FromEmail string
	ApiUrl    string
	SendGrid  sendGridConfig
}
type AppConfig struct {
	Server      serverConfig
	Db          dbConfig
	Mail        MailConfig
	Environment string
	ApiUrl      string
}

func InitializeConfig() (*AppConfig, error) {

	serverAddr := env.GetEnv("ADDR", "8082")
	dbAddr := env.GetEnv("DB_ADDR", "")
	maxOpenConns := env.GetEnvAsInt("DB_MAX_OPEN_CONNS", 30)
	maxIdleConns := env.GetEnvAsInt("DB_MAX_IDLE_CONNS", 30)
	maxIdleTime := env.GetEnvAsDuration("DB_MAX_IDLE_TIME", time.Minute*15)
	apiUrl := env.GetEnv("API_URL", "")
	mailExpiry := env.GetEnvAsDuration("MAIL_EXPIRY", time.Hour*24)
	sendGridAPIKey := env.GetEnv("SENDGRID_API_KEY", "")
	fromEmail := env.GetEnv("FROM_EMAIL", "")
	environment := env.GetEnv("ENV", "development")

	config := &AppConfig{
		Server: serverConfig{
			Host: serverAddr,
		},
		Db: dbConfig{
			Addr:         dbAddr,
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			MaxIdleTime:  maxIdleTime,
		},
		ApiUrl: apiUrl,
		Mail: MailConfig{
			Expiry:    mailExpiry,
			FromEmail: fromEmail,
			ApiUrl:    apiUrl,
			SendGrid: sendGridConfig{
				APIKey: sendGridAPIKey,
			},
		},
		Environment: environment,
	}
	return config, nil
}
