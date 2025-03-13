package main

import (
	"context"
	"fmt"

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
	fmt.Printf("SendNotification event received: %v", in)

	types := make([]string, 0, len(in.Types))
	for _, t := range in.Types {
		types = append(types, t.String())
	}
	
	if len(types) == 0 && len(in.Types) > 0 {
		fmt.Println("Warning: Received notification request with empty type values")
	}

	h.service.Handle(in.TaskId, types)

	return &pb.SendNotificationResponse{
		Ack: "Notification sent",
	}, nil
}
