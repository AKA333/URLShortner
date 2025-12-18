package config

import (
	"os"

	"github.com/AKA333/URLShortner/internal/models"
	"github.com/joho/godotenv"
)

func LoadConfig() *models.Config {
	_ = godotenv.Load()
	
	return &models.Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RedisURL:   getEnv("REDIS_URL", "localhost:6379"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "admin"),
		DBPassword: getEnv("DB_PASSWORD", "securepassword"),
		DBName:     getEnv("DB_NAME", "urlshortener"),
		BaseURL:    getEnv("BASE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}