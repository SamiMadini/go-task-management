package grpc

import (
	"context"

	"sama/go-task-management/commons"

	pb "sama/go-task-management/commons/api"

	"google.golang.org/grpc"
)

type Service struct {
	logger                    commons.Logger
	notificationServiceClient pb.NotificationServiceClient
}

func NewService( logger commons.Logger, notificationServiceClient pb.NotificationServiceClient) *Service {
	return &Service{
		logger:                    logger,
		notificationServiceClient: notificationServiceClient,
	}
}

func (s *Service) SendNotification(ctx context.Context, grpcEvent commons.GRPCEvent) error {
	s.logger.Info("GRPCService::Sending event", grpcEvent)

	notification := &pb.SendNotificationRequest{
		TaskId:        grpcEvent.TaskId,
		CorrelationId: grpcEvent.CorrelationId,
		Types:         convertToNotificationTypes(grpcEvent.Types),
	}

	_, err := s.notificationServiceClient.SendNotification(
		ctx,
		notification,
		grpc.FailFastCallOption{},
	)

	if err != nil {
		s.logger.Error("GRPCService::Failed to send notification", "error", err)
		return err
	}

	return nil
}

func convertToNotificationTypes(types []string) []pb.NotificationType {
	converted := make([]pb.NotificationType, len(types))
	for i, t := range types {
		switch t {
		case "IN_APP":
			converted[i] = pb.NotificationType_IN_APP
		case "EMAIL":
			converted[i] = pb.NotificationType_EMAIL
		case "SMS":
			converted[i] = pb.NotificationType_SMS
		default:
			continue
		}
	}
	return converted
}
