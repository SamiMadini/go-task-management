package main

import (
	"context"
	"log"
	"os"

	"sama/go-task-management/commons"
	"sama/go-task-management/email-service/src/config"
	"sama/go-task-management/email-service/src/handlers"
	"sama/go-task-management/email-service/src/server"
	"sama/go-task-management/email-service/src/sqs"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in main: %v", r)
		}
	}()

	cfg := config.LoadConfig()
	if err := config.InitLogging(cfg); err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}

	log.Println("======= Email Service Starting Up =======")

	database, err := commons.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	eventRepo := commons.NewPostgresTaskSystemEventRepository(database)
	messageHandler := handlers.NewMessageHandler(eventRepo)

	if _, ok := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); ok {
		log.Println("Running in AWS Lambda environment")
		startLambda(messageHandler)
	} else {
		log.Println("Starting local development environment")
		startLocalServer(cfg, messageHandler)
	}
}

func startLambda(handler *handlers.MessageHandler) {
	lambda.Start(func(ctx context.Context, event events.SQSEvent) error {
		log.Printf("Processing Lambda event with %d records", len(event.Records))

		for _, record := range event.Records {
			if err := handler.HandleMessage(ctx, []byte(record.Body)); err != nil {
				log.Printf("Error processing Lambda event: %v", err)
				return err
			}
		}
		return nil
	})
}

func startLocalServer(cfg *config.Config, handler *handlers.MessageHandler) {
	ctx := context.Background()

	sqsManager, err := sqs.NewSQSManager(cfg, handler)
	if err != nil {
		log.Printf("Warning: Failed to initialize SQS manager: %v", err)
		// Continue without SQS in local environment
	}

	srv := server.NewServer(cfg, handler, sqsManager)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
