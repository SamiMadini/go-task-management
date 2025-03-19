package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sama/go-task-management/email-service/src/config"
	"sama/go-task-management/email-service/src/handlers"
	"sama/go-task-management/email-service/src/sqs"
)

type Server struct {
	config        *config.Config
	handler       *handlers.MessageHandler
	sqsManager    *sqs.SQSManager
	httpServer    *http.Server
}

func NewServer(cfg *config.Config, handler *handlers.MessageHandler, sqsManager *sqs.SQSManager) *Server {
	return &Server{
		config:     cfg,
		handler:    handler,
		sqsManager: sqsManager,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.config.Port),
		Handler: s.createHandler(),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("HTTP server listening on port %s", s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	if s.sqsManager != nil {
		s.sqsManager.StartPolling(ctx)
	}

	<-stop

	return s.Shutdown()
}

func (s *Server) Shutdown() error {
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}

func (s *Server) createHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request from %s", r.Method, r.RemoteAddr)

		switch r.Method {
		case http.MethodGet:
			s.handleHealthCheck(w)
		case http.MethodPost:
			s.handleMessage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (s *Server) handleHealthCheck(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email service is running"))
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var payload json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.handler.HandleMessage(r.Context(), payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message processed successfully"))
}
