package memorystorage

import (
	"context"
	"sort"
	"time"

	customerror "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/errors"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

func (s *Storage) GetNotificationsInAdvance(ctx context.Context) ([]models.Notification, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, customerror.CustomError{
			Field:   "",
			Message: ctx.Err().Error(),
		}
	default:
	}

	var notifications []models.Notification

	now := time.Now()
	for _, event := range s.events {
		var notification models.Notification

		if !event.Scheduled && now.Before(event.Date.Add(-event.NotificationInterval)) {
			notification = models.Notification{
				EventID:  event.ID,
				Title:    event.Title,
				Date:     event.Date,
				UserID:   event.UserID,
				Interval: event.NotificationInterval,
			}
		}

		if notification.EventID != "" {
			notifications = append(notifications, notification)
		}
	}

	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].Date.Add(-notifications[i].Interval).Before(notifications[j].Date.Add(-notifications[j].Interval))
	})

	if len(notifications) > 10 {
		notifications = notifications[:10]
	}

	return notifications, nil
}

func (s *Storage) UpdateScheduledNotification(ctx context.Context, id string) error {
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

	event := s.events[id]
	event.Scheduled = true
	s.events[id] = event

	return nil
}
