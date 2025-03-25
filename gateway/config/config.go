package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	HTTPAddress              string
	NotificationServiceAddr  string
	AllowedOrigins          []string
	JWTSecret               string
	JWTExpirationHours      int
	DatabaseURL             string
	Environment             string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		HTTPAddress:             getEnvOrDefault("HTTP_ADDRESS", "localhost:8080"),
		NotificationServiceAddr: getEnvOrDefault("NOTIFICATION_SERVICE_ADDRESS", "localhost:2000"),
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3010",
			"http://localhost:3012",
			"http://localhost:8080",
			"http://frontend:3010",
		},
		JWTSecret:          getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		JWTExpirationHours: getEnvAsIntOrDefault("JWT_EXPIRATION_HOURS", 24),
		DatabaseURL:        getEnvOrDefault("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/task_management?sslmode=disable"),
		Environment:        getEnvOrDefault("ENVIRONMENT", "development"),
	}

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if all required configuration is present
func (c *Config) validate() error {
	if c.JWTSecret == "your-secret-key" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}
	return nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault gets an environment variable as an integer or returns a default value
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
} 