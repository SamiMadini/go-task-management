package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	commons "sama/go-task-management/commons"

	"google.golang.org/grpc"
)

var grpcServerAddr = commons.GetEnv("NOTIFICATION_SERVICE_ADDRESS", "localhost:2000")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	dbConnection, err := commons.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer dbConnection.Close()

	taskRepository := commons.NewPostgresTaskRepository(dbConnection)
	taskSystemEventRepository := commons.NewPostgresTaskSystemEventRepository(dbConnection)
	inAppNotificationRepository := commons.NewPostgresInAppNotificationRepository(dbConnection)

	sqsClient, err := NewSQSClient(ctx, Config{
		AWSEndpoint: commons.GetEnv("AWS_ENDPOINT", ""),
		AWSRegion:   commons.GetEnv("AWS_REGION", defaultRegion),
		QueueName:   commons.GetEnv("QUEUE_NAME", defaultQueueName),
	})
	if err != nil {
		log.Fatalf("Failed to create SQS client: %v", err)
	}

	inAppService := NewInAppNotificationService(taskRepository, taskSystemEventRepository, inAppNotificationRepository)
	emailService := NewEmailNotificationService(taskRepository, taskSystemEventRepository, sqsClient)

	NewGrpcHandler(grpcServer, inAppService, emailService)

	log.Println("Notifications service started at", grpcServerAddr)

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			log.Printf("Failed to serve: %v", err)
			cancel()
		}
	}()

	select {
	case <-sigChan:
		log.Println("Received shutdown signal")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	grpcServer.GracefulStop()
	log.Println("Server stopped gracefully")
}
