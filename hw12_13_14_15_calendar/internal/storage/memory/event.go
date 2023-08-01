package memorystorage

import (
	"context"
	customerror "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/errors"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

const day = 24 * time.Hour

func (s *Storage) CreateEvent(ctx context.Context, event models.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return "", customerror.CustomError{
			Field:   "",
			Message: ctx.Err().Error(),
		}
	default:
	}

	s.events[event.ID] = event

	return event.ID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return models.Event{}, customerror.CustomError{
			Field:   "",
			Message: ctx.Err().Error(),
		}
	default:
	}

	if _, ok := s.events[id]; !ok {
		return models.Event{}, customerror.CustomError{
			Field:   "id",
			Message: "no event with id " + id,
		}
	}

	s.events[id] = event

	return s.events[id], nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return customerror.CustomError{
			Field:   "",
			Message: ctx.Err().Error(),
		}
	default:
	}

	if _, ok := s.events[id]; !ok {
		return customerror.CustomError{
			Field:   "id",
			Message: "no event with id " + id,
		}
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, customerror.CustomError{
			Field:   "",
			Message: ctx.Err().Error(),
		}
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

func (s *Storage) GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, customerror.CustomError{
			Field:   "",
			Message: ctx.Err().Error(),
		}
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

func (s *Storage) GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, customerror.CustomError{
			Field:   "",
			Message: ctx.Err().Error(),
		}
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
