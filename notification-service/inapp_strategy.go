package main

import (
	"context"
	"slices"
)

type InAppNotificationStrategy struct {
	inAppService *InAppNotificationService
}

func NewInAppNotificationStrategy(service *InAppNotificationService) *InAppNotificationStrategy {
	return &InAppNotificationStrategy{inAppService: service}
}

func (s *InAppNotificationStrategy) CanProcess(types []string) bool {
	return slices.Contains(types, "IN_APP")
}

func (s *InAppNotificationStrategy) Process(ctx context.Context, taskID string, correlationID string, types []string) error {
	return s.inAppService.Handle(ctx, taskID, correlationID, types)
}
