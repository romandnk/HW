package storage

import (
	"context"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

type StoreEvent interface {
	Create(ctx context.Context, event models.Event) (string, error)
	Update(ctx context.Context, id string, event models.Event) (models.Event, error)
	Delete(ctx context.Context, id string) (string, error)
	GetAllByDay(ctx context.Context, date time.Time) ([]models.Event, error)
	GetAllByWeek(ctx context.Context, date time.Time) ([]models.Event, error)
	GetAllByMonth(ctx context.Context, date time.Time) ([]models.Event, error)
}
