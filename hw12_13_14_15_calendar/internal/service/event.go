package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

var ErrInvalidUserID = errors.New("user id must be positive number")

func (s *Service) CreateEvent(ctx context.Context, event models.Event) (string, error) {
	if event.UserID <= 0 {
		return "", ErrInvalidUserID
	}
	id := uuid.New().String()
	event.ID = id
	return s.event.CreateEvent(ctx, event)
}

func (s *Service) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	if event.UserID <= 0 {
		return models.Event{}, ErrInvalidUserID
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
