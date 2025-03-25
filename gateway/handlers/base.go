package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	commons "sama/go-task-management/commons"
	pb "sama/go-task-management/commons/api"
	"sama/go-task-management/gateway/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BaseHandler struct {
	userRepository              commons.UserRepositoryInterface
	taskRepository              commons.TaskRepositoryInterface
	taskSystemEventRepository   commons.TaskSystemEventRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
	notificationServiceClient   pb.NotificationServiceClient
	passwordResetTokenRepository commons.PasswordResetTokenRepositoryInterface
	jwtSecret                  string
	logger                     *log.Logger
}

func NewBaseHandler(cfg *config.Config) (*BaseHandler, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		cfg.NotificationServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	database, err := commons.InitDB()
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "[API] ", log.LstdFlags|log.Lshortfile)

	return &BaseHandler{
		userRepository:              commons.NewPostgresUserRepository(database),
		taskRepository:              commons.NewPostgresTaskRepository(database),
		taskSystemEventRepository:   commons.NewPostgresTaskSystemEventRepository(database),
		inAppNotificationRepository: commons.NewPostgresInAppNotificationRepository(database),
		notificationServiceClient:   pb.NewNotificationServiceClient(conn),
		passwordResetTokenRepository: commons.NewPostgresPasswordResetTokenRepository(database),
		jwtSecret:                  os.Getenv("JWT_SECRET"),
		logger:                     logger,
	}, nil
}

func (h *BaseHandler) decodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

func (h *BaseHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response := StandardResponse{
		Success: code >= 200 && code < 300,
		Data:    payload,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		h.logger.Printf("Error marshaling response: %v\nStack trace:\n%s", err, debug.Stack())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func (h *BaseHandler) respondWithError(w http.ResponseWriter, status int, code string, message string, details string) {
	h.respondWithJSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}

func (h *BaseHandler) respondWithValidationErrors(w http.ResponseWriter, errors []ValidationError) {
	h.respondWithJSON(w, http.StatusBadRequest, ErrorResponse{
		Code:             ErrCodeValidation,
		Message:          "Validation failed",
		ValidationErrors: errors,
	})
}

func (h *BaseHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
