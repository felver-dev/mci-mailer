package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SMTP     SMTPConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL string
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	FromName string
	UseTLS   bool
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] no .env file, using environment variables")
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "2525"),
			Env:  getEnv("APP_ENV", "production"),
		},
		Database: DatabaseConfig{
			URL: mustGetEnv("DATABASE_URL"),
		},
		SMTP: SMTPConfig{
			Host:     mustGetEnv("SMTP_HOST"),
			Port:     smtpPort,
			User:     mustGetEnv("SMTP_USER"),
			Password: mustGetEnv("SMTP_PASSWORD"),
			From:     getEnv("SMTP_FROM", getEnv("SMTP_USER", "")),
			FromName: getEnv("SMTP_FROM_NAME", "MCI CARE CI"),
			UseTLS:   getEnv("SMTP_USE_TLS", "true") == "true",
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("[config] required environment variable %q is not set", key)
	}
	return val
}
