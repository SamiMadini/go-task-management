package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                    int
	NotificationServiceAddr string
	AllowedOrigins          []string
	JWTSecret               string
	JWTExpirationHours      int
	DatabaseURL             string
	Environment             string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnvOrDefault("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Port:                    port,
		NotificationServiceAddr: getEnvOrDefault("NOTIFICATION_SERVICE_ADDR", "localhost:50051"),
		AllowedOrigins:          []string{getEnvOrDefault("ALLOWED_ORIGINS", "*")},
		JWTSecret:               getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		JWTExpirationHours:      getEnvAsIntOrDefault("JWT_EXPIRATION_HOURS", 24),
		DatabaseURL:             getEnvOrDefault("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/taskdb?sslmode=disable"),
		Environment:             getEnvOrDefault("ENVIRONMENT", "development"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "your-secret-key" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
