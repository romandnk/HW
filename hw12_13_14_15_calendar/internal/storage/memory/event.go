package memorystorage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"time"
)

const day = 24 * time.Hour

func (s *Storage) Create(event models.Event) string {
	s.mu.Lock()

	id := uuid.New().String()

	event.ID = id

	s.events[id] = event

	s.mu.Unlock()

	return id
}

func (s *Storage) Update(id string, event models.Event) (models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return models.Event{}, fmt.Errorf("updating: no event with id %s", id)
	}

	s.events[id] = event

	return event, nil
}

func (s *Storage) Delete(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return "", fmt.Errorf("deleting: no event with id %s", id)
	}

	delete(s.events, id)

	return id, nil
}

func (s *Storage) GetAllByDay(date time.Time) []models.Event {
	s.mu.RLock()

	var events []models.Event

	for _, event := range s.events {
		if event.Date == date {
			events = append(events, event)
		}
	}

	s.mu.RUnlock()

	return events
}

func (s *Storage) GetAllByWeek(date time.Time) []models.Event {
	s.mu.RLock()

	var events []models.Event

	for _, event := range s.events {
		if inTimeSpan(date, date.Add(6*day), event.Date) {
			events = append(events, event)
		}
	}

	s.mu.RUnlock()

	return events
}

func (s *Storage) GetAllByMonth(date time.Time) []models.Event {
	s.mu.RLock()

	var events []models.Event

	for _, event := range s.events {
		if inTimeSpan(date, date.Add(29*day), event.Date) {
			events = append(events, event)
		}
	}

	s.mu.RUnlock()

	return events
}

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}
