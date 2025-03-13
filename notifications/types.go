package main

import "context"

type NotificationServiceInterface interface {
	Handle(ctx context.Context) error
}
