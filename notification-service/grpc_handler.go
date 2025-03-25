package main

import (
	"context"
	"fmt"
	"log"

	pb "sama/go-task-management/commons/api"

	"google.golang.org/grpc"
)

type handler struct {
	strategies []NotificationStrategy
	pb.UnimplementedNotificationServiceServer
}

func NewGrpcHandler(
	grpcServer *grpc.Server,
	inAppService *InAppNotificationService,
	emailService *EmailNotificationService,
) *handler {
	strategies := []NotificationStrategy{
		NewInAppNotificationStrategy(inAppService),
		NewEmailNotificationStrategy(emailService),
	}
	handler := &handler{strategies: strategies}
	pb.RegisterNotificationServiceServer(grpcServer, handler)
	return handler
}

func (h *handler) validateAndHandleNotificationTypes(ctx context.Context, in *pb.SendNotificationRequest) error {
	types := make([]string, 0, len(in.Types))
	for _, t := range in.Types {
		types = append(types, t.String())
	}

	if len(types) == 0 && len(in.Types) > 0 {
		log.Println("Warning: Received notification request with empty type values")
	}

	processed := false
	for _, strategy := range h.strategies {
		if strategy.CanProcess(types) {
			if err := strategy.Process(ctx, in.TaskId, in.CorrelationId, types); err != nil {
				return fmt.Errorf("failed to process notification: %w", err)
			}
			processed = true
		}
	}

	if !processed {
		return fmt.Errorf("no valid notification strategy found for types: %v", types)
	}

	return nil
}

func (h *handler) SendNotification(ctx context.Context, in *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	if err := h.validateAndHandleNotificationTypes(ctx, in); err != nil {
		return nil, err
	}

	return &pb.SendNotificationResponse{
		Ack: "Notification sent",
	}, nil
}
