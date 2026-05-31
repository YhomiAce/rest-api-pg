package config

import (
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
}

func InitConfig() (*Config, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBName:     getEnv("DB_NAME", "default_db"),
		DBUser:     getEnv("DB_USER", "default_user"),
		DBPassword: getEnv("DB_PASSWORD", "default_password"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}