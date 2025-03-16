package main

import (
	"log"
	"net"

	commons "sama/go-task-management/commons"

	"google.golang.org/grpc"
)

var grpcServerAddr = commons.GetEnv("NOTIFICATION_SERVICE_ADDRESS", "localhost:2000")

func main() {
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

	service := NewNotificationService(taskRepository, taskSystemEventRepository, inAppNotificationRepository)
	
	NewGrpcHandler(grpcServer, service)

	log.Println("Notifications service started at", grpcServerAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	} 
}
