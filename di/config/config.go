package config

import "github.com/kelseyhightower/envconfig"

type AppConfig struct {
	CommonConfig   CommonConfig
	ServerConfig   ServerConfig
	DatabaseConfig DatabaseConfig
	EmailConfig    EmailConfig
}

type CommonConfig struct {
	Env string `envconfig:"ENV" default:"local"`
}

type ServerConfig struct {
	Port         string `envconfig:"APP_PORT" default:"8080"`
	FrontendPath string `envconfig:"APP_FRONTEND_PATH" default:"https://account.tongla.dev"`
	BackendPath  string `envconfig:"APP_BACKEND_PATH" default:"https://account-api.tongla.dev"`
	SentryDns    string `envconfig:"SENTRY_DNS" default:"https://sentry.io"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"3306"`
	User     string `envconfig:"DB_USER" default:"root"`
	Password string `envconfig:"DB_PASSWORD" default:"root"`
	DBName   string `envconfig:"DB_NAME" default:"mydb"`
}

type EmailConfig struct {
	Host     string `envconfig:"SMTP_HOST" default:"localhost@gmail.com"`
	Port     int    `envconfig:"SMTP_PORT" default:"587"`
	Sender   string `envconfig:"SMTP_SENDER" default:"localhost"`
	Password string `envconfig:"SMTP_PASSWORD" default:"localhost"`
}

func GetConfig() AppConfig {
	var app AppConfig
	envconfig.MustProcess("APP", &app.CommonConfig)
	envconfig.MustProcess("APP", &app.ServerConfig)
	envconfig.MustProcess("APP", &app.DatabaseConfig)
	envconfig.MustProcess("APP", &app.EmailConfig)
	return app
}
