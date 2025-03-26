package services

import (
	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/services/adapters"
	"sama/go-task-management/gateway/services/auth"
	"sama/go-task-management/gateway/services/grpc"
	"sama/go-task-management/gateway/services/in_app_notification"
	"sama/go-task-management/gateway/services/task"
	"sama/go-task-management/gateway/services/task_system_event"

	pb "sama/go-task-management/commons/api"
)

type Services struct {
	AuthService              *auth.Service
	TaskService              *task.Service
	TaskSystemEventService   *task_system_event.Service
	InAppNotificationService *in_app_notification.Service
	GrpcService              *grpc.Service
}

func NewServices(
	logger commons.Logger,
	jwtSecret string,
	userRepo commons.UserRepositoryInterface,
	taskRepo commons.TaskRepositoryInterface,
	taskSystemEventRepo commons.TaskSystemEventRepositoryInterface,
	inAppNotificationRepo commons.InAppNotificationRepositoryInterface,
	passwordResetTokenRepo commons.PasswordResetTokenRepositoryInterface,
	notificationServiceClient pb.NotificationServiceClient,
) *Services {
	userAdapter := &adapters.UserRepositoryAdapter{UserRepositoryInterface: userRepo}
	taskAdapter := &adapters.TaskRepositoryAdapter{TaskRepositoryInterface: taskRepo}
	inAppNotificationAdapter := &adapters.InAppNotificationRepositoryAdapter{InAppNotificationRepositoryInterface: inAppNotificationRepo}

	authService := auth.NewService(logger, jwtSecret, userAdapter, passwordResetTokenRepo)
	inAppNotificationService := in_app_notification.NewService(logger, inAppNotificationAdapter)
	taskService := task.NewService(logger, taskAdapter, userAdapter)
	taskSystemEventService := task_system_event.NewService(logger, taskSystemEventRepo)
	grpcService := grpc.NewService(logger, notificationServiceClient)

	return &Services{
		AuthService:              authService,
		TaskService:              taskService,
		TaskSystemEventService:   taskSystemEventService,
		InAppNotificationService: inAppNotificationService,
		GrpcService:              grpcService,
	}
}
