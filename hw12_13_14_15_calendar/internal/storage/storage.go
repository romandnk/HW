package storage

import (
	"context"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

type EventStorage interface {
	CreateEvent(ctx context.Context, event models.Event) (string, error)
	UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	DeleteOutdatedEvents(ctx context.Context) error
	GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error)
	GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error)
	GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error)
}

type NotificationStorage interface {
	UpdateScheduledNotification(ctx context.Context, id string) error
	GetNotificationsInAdvance(ctx context.Context) ([]models.Notification, error)
}

type Storage interface {
	EventStorage
	NotificationStorage
}
