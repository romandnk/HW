package service

import (
	"context"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
)

type NotificationService struct {
	notification storage.NotificationStorage
}

func NewNotificationService(notification storage.NotificationStorage) *NotificationService {
	return &NotificationService{notification: notification}
}

func (n *NotificationService) GetNotificationInAdvance(ctx context.Context) ([]models.Notification, error) {
	return n.notification.GetNotificationInAdvance(ctx)
}
