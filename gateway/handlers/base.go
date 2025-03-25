package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

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

	return &BaseHandler{
		userRepository:              commons.NewPostgresUserRepository(database),
		taskRepository:              commons.NewPostgresTaskRepository(database),
		taskSystemEventRepository:   commons.NewPostgresTaskSystemEventRepository(database),
		inAppNotificationRepository: commons.NewPostgresInAppNotificationRepository(database),
		notificationServiceClient:   pb.NewNotificationServiceClient(conn),
		passwordResetTokenRepository: commons.NewPostgresPasswordResetTokenRepository(database),
		jwtSecret:                  os.Getenv("JWT_SECRET"),
	}, nil
}

func (h *BaseHandler) decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (h *BaseHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *BaseHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *BaseHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
