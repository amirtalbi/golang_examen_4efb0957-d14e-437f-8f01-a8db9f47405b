package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort       string
	JWTSecret        string
	TokenExpiryHours int
	APIPrefix        string
	Database         DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		TokenExpiryHours: getEnvAsInt("TOKEN_EXPIRY_HOURS", 24),
		APIPrefix:        getEnv("API_PREFIX", "567088a9-6689-4e67-b5e5-ed40ad0a830c"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "examen_go"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
