package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	customerror "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/errors"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrInvalidUserID               = errors.New("user id must be positive number")
	ErrEmptyTitle                  = errors.New("title cannot be empty")
	ErrInvalidDuration             = errors.New("duration cannot be non-positive")
	ErrInvalidNotificationInterval = errors.New("notification interval cannot be negative")
)

type EventService struct {
	event storage.EventStorage
}

func NewEventService(event storage.EventStorage) *EventService {
	return &EventService{event: event}
}

func (s *EventService) CreateEvent(ctx context.Context, event models.Event) (string, error) {
	event.Title = strings.TrimSpace(event.Title)
	if event.Title == "" {
		return "", customerror.CustomError{
			Field:   "title",
			Message: ErrEmptyTitle.Error(),
		}
	}
	if event.Duration <= 0 {
		return "", customerror.CustomError{
			Field:   "duration",
			Message: ErrInvalidDuration.Error(),
		}
	}
	event.Description = strings.TrimSpace(event.Description)
	if event.UserID <= 0 {
		return "", customerror.CustomError{
			Field:   "user_id",
			Message: ErrInvalidUserID.Error(),
		}
	}
	if event.NotificationInterval < 0 {
		return "", customerror.CustomError{
			Field:   "notification_interval",
			Message: ErrInvalidNotificationInterval.Error(),
		}
	}

	id := uuid.New().String()
	event.ID = id
	return s.event.CreateEvent(ctx, event)
}

func (s *EventService) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	event.Title = strings.TrimSpace(event.Title)
	if event.Duration < 0 {
		return models.Event{}, customerror.CustomError{
			Field:   "duration",
			Message: ErrInvalidDuration.Error(),
		}
	}
	event.Description = strings.TrimSpace(event.Description)
	if event.UserID < 0 {
		return models.Event{}, customerror.CustomError{
			Field:   "user_id",
			Message: ErrInvalidUserID.Error(),
		}
	}
	if event.NotificationInterval < 0 {
		return models.Event{}, customerror.CustomError{
			Field:   "notification_interval",
			Message: ErrInvalidNotificationInterval.Error(),
		}
	}

	return s.event.UpdateEvent(ctx, id, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, id string) error {
	return s.event.DeleteEvent(ctx, id)
}

func (s *EventService) GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.event.GetAllByDayEvents(ctx, date)
}

func (s *EventService) GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.event.GetAllByWeekEvents(ctx, date)
}

func (s *EventService) GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.event.GetAllByMonthEvents(ctx, date)
}
