package memorystorage

import (
	"context"
	"time"

	customerror "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/errors"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

func (s *Storage) GetNotificationInAdvance(ctx context.Context) ([]models.Notification, error) {
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
				EventID: event.ID,
				Title:   event.Title,
				Date:    event.Date,
				UserID:  event.UserID,
			}
		}

		if notification.EventID != "" {
			event.Scheduled = true
			s.events[event.ID] = event
			notifications = append(notifications, notification)
		}
	}

	return notifications, nil
}
