package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SMTP     SMTPConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port         string
	Env          string
	CORSOrigins  []string
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

type AuthConfig struct {
	JWTSecret   string
	MasterToken string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] no .env file, using environment variables")
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	rawOrigins := getEnv("CORS_ORIGINS", "http://localhost:5173")
	origins := strings.Split(rawOrigins, ",")

	return &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "2525"),
			Env:         getEnv("APP_ENV", "production"),
			CORSOrigins: origins,
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
		Auth: AuthConfig{
			JWTSecret:   mustGetEnv("JWT_SECRET"),
			MasterToken: mustGetEnv("MASTER_TOKEN"),
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
