package postgres

import (
	"context"
	"fmt"
	"time"

	customerror "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/errors"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

func (s *Storage) GetNotificationsInAdvance(ctx context.Context) ([]models.Notification, error) {
	var notifications []models.Notification

	selectNotifications := fmt.Sprintf(`
		SELECT id, title, date, user_id, notification_interval
		FROM %s
		WHERE scheduled = false AND $1::timestamp <= date - notification_interval
		ORDER BY (date - notification_interval)
		LIMIT 10;`, eventsTable)

	nw := time.Now().UTC()

	rows, err := s.db.Query(ctx, selectNotifications, nw)
	if err != nil {
		return nil, customerror.CustomError{
			Field:   "",
			Message: err.Error(),
		}
	}

	for rows.Next() {
		var notification models.Notification

		err = rows.Scan(
			&notification.EventID,
			&notification.Title,
			&notification.Date,
			&notification.UserID,
			&notification.Interval)
		if err != nil {
			return nil, customerror.CustomError{
				Field:   "",
				Message: err.Error(),
			}
		}

		notifications = append(notifications, notification)
	}

	if rows.Err() != nil {
		return nil, customerror.CustomError{
			Field:   "",
			Message: err.Error(),
		}
	}

	return notifications, nil
}

func (s *Storage) UpdateScheduledNotification(ctx context.Context, id string) error {
	updateNotifications := fmt.Sprintf(`
		UPDATE %s
		SET scheduled = true 
		WHERE id = $1`, eventsTable)

	ct, err := s.db.Exec(ctx, updateNotifications, id)
	if err != nil {
		return customerror.CustomError{
			Field:   "",
			Message: err.Error(),
		}
	}

	if ct.RowsAffected() == 0 {
		return customerror.CustomError{
			Field:   "id",
			Message: "notification wasn't updated with id: " + id,
		}
	}

	return nil
}
