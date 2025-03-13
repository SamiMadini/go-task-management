package main

import (
	"log"
	"net/http"

	commons "sama/go-task-management/commons"

	pb "sama/go-task-management/commons/api"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var httpAddress = commons.GetEnv("HTTP_ADDRESS", "localhost:8080")
var notificationServiceAddress = commons.GetEnv("NOTIFICATION_SERVICE_ADDRESS", "localhost:2000")

func main() {
	conn, err := grpc.Dial(notificationServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to notification service at %s", notificationServiceAddress)

	notificationServiceClient := pb.NewNotificationServiceClient(conn)

	database, err := commons.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Printf("Connected to database")

	taskRepository := commons.NewPostgresTaskRepository(database)
	inAppNotificationRepository := commons.NewPostgresInAppNotificationRepository(database)

	mux := http.NewServeMux()

	handler := NewHandler(taskRepository, inAppNotificationRepository, notificationServiceClient)
	handler.registerRoutes(mux)

	log.Printf("Starting server on %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
