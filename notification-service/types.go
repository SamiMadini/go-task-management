package main

import "context"

type NotificationServiceInterface interface {
	Handle(ctx context.Context) error
}

type NotificationStrategy interface {
	CanProcess(types []string) bool
	Process(ctx context.Context, taskID string, correlationID string, types []string) error
}
