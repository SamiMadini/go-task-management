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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sama/go-task-management/gateway/config"
	"sama/go-task-management/gateway/handlers"
	"sama/go-task-management/gateway/middleware"
	"sama/go-task-management/gateway/routes"

	_ "sama/go-task-management/gateway/docs"
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
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create router
	router := routes.NewRouter()

	// Create handler instance
	handler, err := handlers.NewHandler(cfg)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	// Register routes
	router.RegisterRoutes(handler, middleware.DefaultCorsConfig(), middleware.DefaultAuthConfig(cfg.JWTSecret))

	// Create middleware chain
	chain := middleware.NewChain(
		middleware.RecoveryMiddleware,
		middleware.LoggingMiddleware,
		middleware.CorsMiddleware(middleware.DefaultCorsConfig()),
		middleware.AuthMiddleware(middleware.DefaultAuthConfig(cfg.JWTSecret)),
	)

	// Create server with middleware chain
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: chain.Then(router.GetMux()),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
