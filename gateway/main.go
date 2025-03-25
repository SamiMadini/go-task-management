package main

// @title Task Management API
// @version 1.0
// @description API for managing tasks, notifications, and system events
// @host localhost:3012
// @BasePath /api/v1

import (
	"context"
	"log"
	"net/http"

	commons "sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"
	"sama/go-task-management/gateway/config"
	"sama/go-task-management/gateway/middleware"

	_ "sama/go-task-management/gateway/docs"

	_ "github.com/joho/godotenv/autoload"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ErrorResponse represents an error response
// @Description Error response model
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid task ID format"`
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		cfg.NotificationServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to notification service at %s", cfg.NotificationServiceAddr)

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

	handler.registerRoutes(mux)

	chain := middleware.NewChain(
		middleware.RecoveryMiddleware,
		middleware.LoggingMiddleware,
		middleware.CorsMiddleware(middleware.CorsConfig{
			AllowedOrigins: cfg.AllowedOrigins,
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With", "Accept"},
			MaxAge:         86400,
		}),
		middleware.AuthMiddleware(middleware.DefaultAuthConfig(cfg.JWTSecret)),
	)

	log.Printf("Starting server on %s", cfg.HTTPAddress)
	if err := http.ListenAndServe(cfg.HTTPAddress, chain.Then(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
