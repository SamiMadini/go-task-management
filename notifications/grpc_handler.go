package main

import (
	"context"
	"log"

	pb "sama/go-task-management/commons/api"

	"google.golang.org/grpc"
)

type handler struct {
	service *NotificationService
	pb.UnimplementedNotificationServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server, service *NotificationService) *handler {
	handler := &handler{service: service}
	pb.RegisterNotificationServiceServer(grpcServer, handler)
	return handler
}

func (h *handler) SendNotification(ctx context.Context, in *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	types := make([]string, 0, len(in.Types))
	for _, t := range in.Types {
		types = append(types, t.String())
	}
	
	if len(types) == 0 && len(in.Types) > 0 {
		log.Println("Warning: Received notification request with empty type values")
	}

	h.service.Handle(in.TaskId, in.CorrelationId, types)

	return &pb.SendNotificationResponse{
		Ack: "Notification sent",
	}, nil
}
