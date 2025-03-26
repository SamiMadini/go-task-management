package main

// @title Task Management API
// @version 1.0
// @description API for managing tasks, notifications, and system events
// @description
// @description Error Handling:
// @description - 400 Bad Request: Invalid input, validation errors, or malformed requests
// @description - 401 Unauthorized: Missing or invalid authentication token
// @description - 403 Forbidden: Valid token but insufficient permissions
// @description - 404 Not Found: Resource not found
// @description - 500 Internal Server Error: Unexpected server errors
// @host localhost:3012
// @BasePath /api/v1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/config"
	"sama/go-task-management/gateway/handlers"
	"sama/go-task-management/gateway/middleware"
	"sama/go-task-management/gateway/routes"
	"sama/go-task-management/gateway/services"

	_ "sama/go-task-management/gateway/docs"

	pb "sama/go-task-management/commons/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ErrorResponse represents a standardized error response
// @Description Error response model with detailed information about the error
// @Description Common error codes:
// @Description - BAD_REQUEST: Invalid input or request format
// @Description - UNAUTHORIZED: Authentication required
// @Description - FORBIDDEN: Permission denied
// @Description - NOT_FOUND: Resource not found
// @Description - INTERNAL_ERROR: Server error
// @Description - VALIDATION_ERROR: Input validation failed
type ErrorResponse struct {
	Code             string            `json:"code" example:"BAD_REQUEST" enums:"BAD_REQUEST,UNAUTHORIZED,FORBIDDEN,NOT_FOUND,INTERNAL_ERROR,VALIDATION_ERROR"`
	Message          string            `json:"message" example:"Invalid task format"`
	Details          string            `json:"details,omitempty" example:"Task ID must be a valid UUID"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// ValidationError represents a field-level validation error
// @Description Validation error for a specific field
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email address is invalid"`
}

func main() {
	// Initialize logger
	logger := commons.NewLogger("[API] ")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Error("Warning: Error loading .env file:", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration:", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := commons.InitDB()
	if err != nil {
		logger.Error("Failed to initialize database:", err)
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := commons.NewPostgresUserRepository(db)
	taskRepo := commons.NewPostgresTaskRepository(db)
	taskSystemEventRepo := commons.NewPostgresTaskSystemEventRepository(db)
	inAppNotificationRepo := commons.NewPostgresInAppNotificationRepository(db)
	passwordResetTokenRepo := commons.NewPostgresPasswordResetTokenRepository(db)

	// Initialize GRPC service client
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		cfg.NotificationServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("Failed to initialize notification service:", err)
		// os.Exit(1)
	}
	notificationServiceClient := pb.NewNotificationServiceClient(conn)

	// Initialize services
	services := services.NewServices(
		logger,
		os.Getenv("JWT_SECRET"),
		userRepo,
		taskRepo,
		taskSystemEventRepo,
		inAppNotificationRepo,
		passwordResetTokenRepo,
		notificationServiceClient,
	)

	// Initialize handlers
	h, err := handlers.NewHandlers(logger, services)
	if err != nil {
		logger.Error("Failed to initialize handlers:", err)
		os.Exit(1)
	}

	// Create handler wrapper
	handlerWrapper := handlers.NewHandlerWrapper(
		h.Base,
		h.Auth,
		h.Task,
		h.InAppNotification,
		h.TaskSystemEvent,
	)

	// Initialize router
	router := routes.NewRouter()

	// Create base middleware chain
	baseChain := middleware.NewChain(
		middleware.CorsMiddleware(middleware.DefaultCorsConfig()),
		chiMiddleware.Logger,
		chiMiddleware.Recoverer,
		chiMiddleware.RequestID,
		chiMiddleware.RealIP,
		chiMiddleware.Timeout(60*time.Second),
	)

	// Get Chi router and apply base middleware
	r := router.GetRouter()
	r.Use(baseChain.Then)

	// Register routes
	authConfig := middleware.DefaultAuthConfig(os.Getenv("JWT_SECRET"))
	router.RegisterRoutes(handlerWrapper, authConfig)

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		logger.Infof("Server started on port %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server:", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown:", err)
		os.Exit(1)
	}

	logger.Info("Server exited properly")
}
