package postgres

import (
	"context"
	"fmt"
	"time"

	customerror "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/errors"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

func (s *Storage) GetNotificationInAdvance(ctx context.Context) ([]models.Notification, error) {
	var notifications []models.Notification

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, customerror.CustomError{
			Field:   "",
			Message: err.Error(),
		}
	}
	defer tx.Rollback(ctx)

	selectNotifications := fmt.Sprintf(`
		SELECT id, title, date, user_id
		FROM %s
		WHERE scheduled = false AND $1 <= date - notification_interval
		ORDER BY date DESC
		LIMIT 10;`, eventsTable)

	nw := time.Now().Format(time.RFC3339Nano)

	rows, err := tx.Query(ctx, selectNotifications, nw)
	if err != nil {
		return nil, customerror.CustomError{
			Field:   "",
			Message: err.Error(),
		}
	}

	for rows.Next() {
		var notification models.Notification

		err = rows.Scan(&notification.EventID, &notification.Title, &notification.Date, &notification.UserID)
		if err != nil {
			return nil, customerror.CustomError{
				Field:   "",
				Message: err.Error(),
			}
		}

		updateEvents := fmt.Sprintf(`
			UPDATE %s
			SET scheduled = true
			WHERE id = $1`, eventsTable)

		_, err := tx.Exec(ctx, updateEvents, notification.EventID)
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

	err = tx.Commit(ctx)
	if err != nil {
		return nil, customerror.CustomError{
			Field:   "",
			Message: err.Error(),
		}
	}

	return notifications, nil
}
