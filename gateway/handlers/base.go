package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	commons "sama/go-task-management/commons"
	"sama/go-task-management/gateway/config"
	"sama/go-task-management/gateway/handlers/constants"
	"sama/go-task-management/gateway/handlers/validation"
	"sama/go-task-management/gateway/services/auth"
	"sama/go-task-management/gateway/services/grpc"
	in_app_notification "sama/go-task-management/gateway/services/in_app_notification"
	"sama/go-task-management/gateway/services/task"
	"sama/go-task-management/gateway/services/task_system_event"
)

type BaseHandler struct {
	config              *Config
	jwtSecret           string
	logger              commons.Logger
	grpcService         *grpc.Service
	authService         *auth.Service
	taskService         *task.Service
	taskEventService    *task_system_event.Service
	inAppNotificationService *in_app_notification.Service
}

type Config struct {
	JWTSecret string
}

func NewBaseHandler(
	cfg *config.Config,
	logger commons.Logger,
	config *Config,
	authService *auth.Service,
	taskService *task.Service,
	taskEventService *task_system_event.Service,
	inAppNotificationService *in_app_notification.Service,
	grpcService *grpc.Service,
) (*BaseHandler, error) {
	return &BaseHandler{
		config:              config,
		jwtSecret:           os.Getenv("JWT_SECRET"),
		logger:              logger,
		grpcService:         grpcService,
		authService:         authService,
		taskService:         taskService,
		taskEventService:    taskEventService,
		inAppNotificationService: inAppNotificationService,
	}, nil
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (h *BaseHandler) respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("Error encoding response: %v", err)
	}
}

func (h *BaseHandler) respondWithError(w http.ResponseWriter, status int, code string, message string, details string) {
	response := StandardResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	h.respondWithJSON(w, status, response)
}

func (h *BaseHandler) respondWithValidationErrors(w http.ResponseWriter, errors []validation.ValidationError) {
	response := StandardResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    constants.ErrCodeValidation,
			Message: "Validation failed",
			Details: "One or more fields failed validation",
		},
		Data: errors,
	}
	h.respondWithJSON(w, http.StatusBadRequest, response)
}

func (h *BaseHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
