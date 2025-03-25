package main

import (
	"context"
	"slices"
)

type EmailNotificationStrategy struct {
	emailService *EmailNotificationService
}

func NewEmailNotificationStrategy(service *EmailNotificationService) *EmailNotificationStrategy {
	return &EmailNotificationStrategy{emailService: service}
}

func (s *EmailNotificationStrategy) CanProcess(types []string) bool {
	return slices.Contains(types, "EMAIL")
}

func (s *EmailNotificationStrategy) Process(ctx context.Context, taskID string, correlationID string, types []string) error {
	return s.emailService.Handle(ctx, taskID, correlationID, types)
}
