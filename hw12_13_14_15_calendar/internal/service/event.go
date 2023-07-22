package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

func (s *Service) CreateEvent(ctx context.Context, event models.Event) (string, error) {
	id := uuid.New().String()
	event.ID = id
	return s.storage.CreateEvent(ctx, event) //nolint:nolintlint,gci
}

func (s *Service) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	return s.storage.UpdateEvent(ctx, id, event) //nolint:nolintlint,gci
}

func (s *Service) DeleteEvent(ctx context.Context, id string) error {
	return s.storage.DeleteEvent(ctx, id) //nolint:nolintlint,gci
}

func (s *Service) GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.storage.GetAllByDayEvents(ctx, date) //nolint:nolintlint,gci
}

func (s *Service) GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.storage.GetAllByWeekEvents(ctx, date) //nolint:nolintlint,gci
}

func (s *Service) GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	return s.storage.GetAllByMonthEvents(ctx, date)
}
