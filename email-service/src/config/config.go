package config

import (
	"io"
	"log"
	"os"
)

type Config struct {
	Port           string
	QueueName      string
	AWSEndpoint    string
	AWSRegion      string
	LogFilePath    string
	PostgresConfig PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func InitLogging(config *Config) error {
	logFile, err := os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	return nil
}

func LoadConfig() *Config {
	config := &Config{
		Port:        getEnvOrDefault("PORT", "8080"),
		QueueName:   getEnvOrDefault("SQS_QUEUE_NAME", "go-email-service-queue"),
		AWSEndpoint: os.Getenv("AWS_ENDPOINT_URL"),
		AWSRegion:   getEnvOrDefault("AWS_REGION", "us-east-1"),
		LogFilePath: "/tmp/email-service.log",
		PostgresConfig: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
		},
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
