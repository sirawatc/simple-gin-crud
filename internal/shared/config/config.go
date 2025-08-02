package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string
	Database    DatabaseConfig
	Server      ServerConfig
}

type DatabaseConfig struct {
	User        string
	Password    string
	Host        string
	Port        string
	DBName      string
	SSLMode     string
	TimeZone    string
	AutoMigrate bool
}

type ServerConfig struct {
	Host string
	Port string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using default values")
	}
	return &Config{
		ServiceName: getValue("SERVICE_NAME", "simple-gin-crud"),
		Database: DatabaseConfig{
			User:        getValue("DB_USER", ""),
			Password:    getValue("DB_PASSWORD", ""),
			Host:        getValue("DB_HOST", ""),
			Port:        getValue("DB_PORT", ""),
			DBName:      getValue("DB_NAME", ""),
			SSLMode:     getValue("DB_SSLMODE", ""),
			TimeZone:    getValue("DB_TIMEZONE", ""),
			AutoMigrate: getValue("DB_AUTO_MIGRATE", "false") == "true",
		},
		Server: ServerConfig{
			Host: getValue("SERVER_HOST", "0.0.0.0"),
			Port: getValue("SERVER_PORT", "8080"),
		},
	}
}

func getValue(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
