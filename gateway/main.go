package main

// @title Task Management API
// @version 1.0
// @description API for managing tasks, notifications, and system events
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
	"sama/go-task-management/gateway/middleware"
	"sama/go-task-management/gateway/routes"

	_ "sama/go-task-management/gateway/docs"
)

// ErrorResponse represents an error response
// @Description Error response model
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid task ID format"`
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
	handler := NewHandler(cfg)

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
		Addr:    fmt.Sprintf(":%d", cfg.Port),
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
