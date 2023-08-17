package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
)

type Event interface {
	CreateEvent(ctx context.Context, event models.Event) (string, error)
	UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	DeleteOutdatedEvents(ctx context.Context) error
	GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error)
	GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error)
	GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error)
}

type Notification interface {
	GetNotificationInAdvance(ctx context.Context) ([]models.Notification, error)
}

type Services interface {
	Event
	Notification
}

type Service struct {
	Event
	Notification
}

func NewService(repo storage.Storage) *Service {
	return &Service{
		NewEventService(repo),
		NewNotificationService(repo),
	}
}
