package memorystorage

import (
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"sync"
	"time"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]models.Event
}

func NewStorage() *Storage {
	return &Storage{
		events: make(map[string]models.Event),
	}
}

type StoreEvent interface {
	Create(event models.Event) string
	Update(id string, event models.Event) (models.Event, error)
	Delete(id string) (string, error)
	GetAllByDay(date time.Time) []models.Event
	GetAllByWeek(date time.Time) []models.Event
	GetAllByMonth(date time.Time) []models.Event
}
