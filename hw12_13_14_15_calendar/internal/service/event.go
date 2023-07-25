package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

var (
	ErrInvalidUserID               = errors.New("user id must not be negative number")
	ErrEmptyTitle                  = errors.New("title cannot be empty")
	ErrInvalidDuration             = errors.New("duration cannot be non-positive")
	ErrInvalidNotificationInterval = errors.New("notification interval cannot be negative")
)

func (s *Service) CreateEvent(ctx context.Context, event models.Event) (string, error) {
	if event.Title == "" {
		return "", ErrEmptyTitle
	}
	if event.Duration <= 0 {
		return "", ErrInvalidDuration
	}
	if event.UserID <= 0 {
		return "", ErrInvalidUserID
	}
	if event.NotificationInterval < 0 {
		return "", ErrInvalidNotificationInterval
	}

	id := uuid.New().String()
	event.ID = id
	return s.event.CreateEvent(ctx, event)
}

func (s *Service) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	var emptyEvent models.Event
	if event.Duration < 0 {
		return emptyEvent, ErrInvalidDuration
	}
	if event.UserID < 0 {
		return emptyEvent, ErrInvalidUserID
	}
	if event.NotificationInterval < 0 {
		return emptyEvent, ErrInvalidNotificationInterval
	}

	return s.event.UpdateEvent(ctx, id, event)
}

func (s *Service) DeleteEvent(ctx context.Context, id string) error {
	return s.event.DeleteEvent(ctx, id)
}

func (s *Service) GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.event.GetAllByDayEvents(ctx, date)
}

func (s *Service) GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.event.GetAllByWeekEvents(ctx, date)
}

func (s *Service) GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.event.GetAllByMonthEvents(ctx, date)
}
