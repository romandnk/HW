package memorystorage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

const day = 24 * time.Hour

func (s *Storage) Create(ctx context.Context, event models.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	id := uuid.New().String()

	event.ID = id

	s.events[id] = event

	return id, nil
}

func (s *Storage) Update(ctx context.Context, id string, event models.Event) (models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return models.Event{}, ctx.Err()
	default:
	}

	if _, ok := s.events[id]; !ok {
		return models.Event{}, fmt.Errorf("updating: no event with id %s", id)
	}

	s.events[id] = event

	return s.events[id], nil
}

func (s *Storage) Delete(ctx context.Context, id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if _, ok := s.events[id]; !ok {
		return "", fmt.Errorf("deleting: no event with id %s", id)
	}

	delete(s.events, id)

	return id, nil
}

func (s *Storage) GetAllByDay(ctx context.Context, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return []models.Event{}, ctx.Err()
	default:
	}

	var events []models.Event

	for _, event := range s.events {
		if event.Date == date {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) GetAllByWeek(ctx context.Context, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return []models.Event{}, ctx.Err()
	default:
	}

	var events []models.Event

	for _, event := range s.events {
		if inTimeSpan(date, date.Add(6*day), event.Date) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) GetAllByMonth(ctx context.Context, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return []models.Event{}, ctx.Err()
	default:
	}

	var events []models.Event

	for _, event := range s.events {
		if inTimeSpan(date, date.Add(29*day), event.Date) {
			events = append(events, event)
		}
	}

	return events, nil
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
