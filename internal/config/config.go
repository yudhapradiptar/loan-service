package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Notification NotificationConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type NotificationConfig struct {
	BaseURL string
	APIKey  string
}

func LoadEnv() error {
	return godotenv.Load()
}

func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "loan_service"),
		},
		Notification: NotificationConfig{
			BaseURL: getEnv("NOTIFICATION_BASE_URL", "http://localhost:8080"),
			APIKey:  getEnv("NOTIFICATION_API_KEY", "1234567890"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
