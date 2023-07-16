package memorystorage

import (
	"sync"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]models.Event
}

func NewStorageMemory() *Storage {
	return &Storage{
		events: make(map[string]models.Event),
	}
}
