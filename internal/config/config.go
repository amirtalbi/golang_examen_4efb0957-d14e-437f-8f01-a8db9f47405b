package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort       string
	JWTSecret        string
	ResetTokenSecret string
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
		ResetTokenSecret: getEnv("RESET_TOKEN_SECRET", "reset-token-secret-key"),
		TokenExpiryHours: getEnvAsInt("TOKEN_EXPIRY_HOURS", 24),
		APIPrefix:        getEnv("API_PREFIX", "4efb0957-d14e-437f-8f01-a8db9f47405b"),
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
