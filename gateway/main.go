package main

// @title Task Management API
// @version 1.0
// @description API for managing tasks, notifications, and system events
// @host localhost:3012
// @BasePath /api/v1

import (
	"log"
	"net/http"

	commons "sama/go-task-management/commons"

	pb "sama/go-task-management/commons/api"

	_ "sama/go-task-management/gateway/docs"

	_ "github.com/joho/godotenv/autoload"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
)

// ErrorResponse represents an error response
// @Description Error response model
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid task ID format"`
}

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

	userRepository := commons.NewPostgresUserRepository(database)
	taskRepository := commons.NewPostgresTaskRepository(database)
	taskSystemEventRepository := commons.NewPostgresTaskSystemEventRepository(database)
	inAppNotificationRepository := commons.NewPostgresInAppNotificationRepository(database)
	passwordResetTokenRepository := commons.NewPostgresPasswordResetTokenRepository(database)

	handler := NewHandler(
		userRepository,
		taskRepository,
		taskSystemEventRepository,
		inAppNotificationRepository,
		notificationServiceClient,
		passwordResetTokenRepository,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:3012/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
		).ServeHTTP(w, r)
	})

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowedOrigins := []string{"http://localhost:3000", "http://localhost:3010", "http://localhost:3012"}
			
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

	handler.registerRoutes(mux)

	log.Printf("Starting server on %s", httpAddress)

	chainMiddlewares := corsMiddleware(AuthMiddleware(mux))

	if err := http.ListenAndServe(httpAddress, chainMiddlewares); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
