package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type EmailNotificationEvent struct {
	TaskId          string  `json:"taskId"`
	CorrelationId   string  `json:"correlationId"`
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	var emailNotificationEvent EmailNotificationEvent
	if err := json.Unmarshal(event, &emailNotificationEvent); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return err
	}

	log.Printf("Email service")

	log.Printf("PostgreSQL Configuration: Host=%s, Port=%s, User=%s, DB=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"))

	log.Printf("Request parameters: taskId=%s, correlationId=%s",
		emailNotificationEvent.TaskId, emailNotificationEvent.CorrelationId)

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
