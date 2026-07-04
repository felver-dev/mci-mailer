package config

import (
	"fmt"
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
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	URL      string
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
		Database: buildDatabaseConfig(),
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

func buildDatabaseConfig() DatabaseConfig {
	host     := getEnv("DB_HOST", "localhost")
	port     := getEnv("DB_PORT", "5432")
	user     := mustGetEnv("DB_USER")
	password := mustGetEnv("DB_PASSWORD")
	name     := mustGetEnv("DB_NAME")
	sslMode  := getEnv("DB_SSLMODE", "disable")

	url := getEnv("DATABASE_URL", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, name, sslMode,
	))

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		SSLMode:  sslMode,
		URL:      url,
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
