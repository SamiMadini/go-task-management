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

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowedOrigins := []string{"http://localhost:3000", "http://localhost:3010"}
			
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			if w.Header().Get("Access-Control-Allow-Origin") == "" {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}

	handler := NewHandler(taskRepository, inAppNotificationRepository, notificationServiceClient)
	handler.registerRoutes(mux)

	log.Printf("Starting server on %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, corsMiddleware(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
