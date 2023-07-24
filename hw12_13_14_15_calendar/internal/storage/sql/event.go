package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

func (s *Storage) CreateEvent(ctx context.Context, event models.Event) (string, error) {
	var id string

	query := fmt.Sprintf(`
		INSERT INTO %s (id, title, date, duration, description, user_id, notification_interval)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`, eventsTable)

	err := s.db.QueryRow(ctx, query,
		event.ID,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationInterval).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	var updatedEvent models.Event

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return updatedEvent, err
	}
	defer func(err *error) {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			*err = rbErr
		}
	}(&err)

	querySelect := fmt.Sprintf(`
		SELECT title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE id = $1`, eventsTable)

	err = tx.QueryRow(ctx, querySelect, id).Scan(
		&updatedEvent.Title,
		&updatedEvent.Date,
		&updatedEvent.Duration,
		&updatedEvent.Description,
		&updatedEvent.UserID,
		&updatedEvent.NotificationInterval,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return updatedEvent, fmt.Errorf("no event with id %s", id)
		}
		return updatedEvent, err
	}

	if event.Title != "" {
		updatedEvent.Title = event.Title
	}
	if event.Date.IsZero() {
		updatedEvent.Date = event.Date
	}
	if event.Duration != 0 {
		updatedEvent.Duration = event.Duration
	}
	if event.Description != "" {
		updatedEvent.Description = event.Description
	}
	if event.UserID != 0 {
		updatedEvent.UserID = event.UserID
	}
	if event.NotificationInterval != 0 {
		updatedEvent.NotificationInterval = event.NotificationInterval
	}

	queryUpdate := fmt.Sprintf(`
		UPDATE %s SET 
        	title = $1, 
            date = $2, 
            duration = $3, 
            description = $4, 
            user_id = $5, 
            notification_interval = $6 
        WHERE id = $7 
        RETURNING id, title, date, duration, description, user_id, notification_interval`, eventsTable)

	err = s.db.QueryRow(ctx, queryUpdate,
		updatedEvent.Title,
		updatedEvent.Date,
		updatedEvent.Duration,
		updatedEvent.Description,
		updatedEvent.UserID,
		updatedEvent.NotificationInterval,
		id,
	).Scan(&updatedEvent.ID, &updatedEvent.Title, &updatedEvent.Date, &updatedEvent.Duration,
		&updatedEvent.Description, &updatedEvent.UserID, &updatedEvent.NotificationInterval)
	if err != nil {
		return updatedEvent, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return updatedEvent, err
	}

	return updatedEvent, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, eventsTable)

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("no event with id %s", id)
	}

	return nil
}

func (s *Storage) GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event

	query := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date = $1`, eventsTable)

	rows, err := s.db.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Date,
			&event.Duration,
			&event.Description,
			&event.UserID,
			&event.NotificationInterval,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event

	query := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date BETWEEN $1 AND $1 + INTERVAL '7 days'`, eventsTable)

	rows, err := s.db.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Date,
			&event.Duration,
			&event.Description,
			&event.UserID,
			&event.NotificationInterval,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event

	query := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date BETWEEN $1 AND $1 + INTERVAL '1 month'`, eventsTable)

	rows, err := s.db.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Date,
			&event.Duration,
			&event.Description,
			&event.UserID,
			&event.NotificationInterval,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
