package handlers

import (
	"os"
	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/services"
)

type Handlers struct {
	Base              *BaseHandler
	Auth              *AuthHandler
	Task              *TaskHandler
	InAppNotification *InAppNotificationHandler
	TaskSystemEvent   *TaskSystemEventHandler
}

func NewHandlers(logger commons.Logger, services *services.Services) (*Handlers, error) {
	baseHandler, err := NewBaseHandler(
		nil, // config will be set later if needed
		logger,
		&Config{
			JWTSecret: os.Getenv("JWT_SECRET"),
		},
		services.AuthService,
		services.TaskService,
		services.TaskSystemEventService,
		services.InAppNotificationService,
		services.GrpcService,
	)
	if err != nil {
		return nil, err
	}

	return &Handlers{
		Base:              baseHandler,
		Auth:              NewAuthHandler(baseHandler, services.AuthService),
		Task:              NewTaskHandler(baseHandler, services.TaskService),
		TaskSystemEvent:   NewTaskSystemEventHandler(baseHandler, services.TaskSystemEventService),
		InAppNotification: NewInAppNotificationHandler(baseHandler, services.InAppNotificationService),
	}, nil
}

func NewHandlerWrapper(
	base *BaseHandler,
	auth *AuthHandler,
	task *TaskHandler,
	inApp *InAppNotificationHandler,
	taskSystem *TaskSystemEventHandler,
	) *HandlerWrapper {
	return &HandlerWrapper{
		Base:              base,
		Auth:              auth,
		Task:              task,
		InAppNotification: inApp,
		TaskSystemEvent:   taskSystem,
	}
}
