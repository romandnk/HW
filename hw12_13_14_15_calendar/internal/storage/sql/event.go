package sqlstorage

import (
	"context"
	"database/sql"
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

type eventForUpdate struct {
	Title                sql.NullString
	Date                 sql.NullTime
	Duration             sql.NullString
	Description          sql.NullString
	UserID               sql.NullInt64
	NotificationInterval sql.NullString
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	var updatedEvent models.Event

	query := fmt.Sprintf(`
		UPDATE %s SET 
        	title = COALESCE($1, title),
            date = COALESCE($2, date),
            duration = COALESCE($3, duration),
            description = COALESCE($4, description),
            user_id = COALESCE($5, user_id),
            notification_interval = COALESCE($6, notification_interval)
        WHERE id = $7 
        RETURNING id, title, date, duration, description, user_id, notification_interval`, eventsTable)

	eventUpdating := checkEmptyFields(event)

	err := s.db.QueryRow(ctx, query,
		eventUpdating.Title,
		eventUpdating.Date,
		eventUpdating.Duration,
		eventUpdating.Description,
		eventUpdating.UserID,
		eventUpdating.NotificationInterval,
		id).Scan(&updatedEvent.ID, &updatedEvent.Title, &updatedEvent.Date, &updatedEvent.Duration,
		&updatedEvent.Description, &updatedEvent.UserID, &updatedEvent.NotificationInterval)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return updatedEvent, fmt.Errorf("no event with id %s", id)
		}
		return updatedEvent, err
	}

	return updatedEvent, nil
}

func checkEmptyFields(event models.Event) eventForUpdate {
	return eventForUpdate{
		Title: sql.NullString{
			String: event.Title,
			Valid:  event.Title != "",
		},
		Date: sql.NullTime{
			Time:  event.Date,
			Valid: !event.Date.IsZero(),
		},
		Duration: sql.NullString{
			String: event.Duration.String(),
			Valid:  event.Duration != 0,
		},
		Description: sql.NullString{
			String: event.Description,
			Valid:  event.Description != "",
		},
		UserID: sql.NullInt64{
			Int64: int64(event.UserID),
			Valid: event.UserID != 0,
		},
		NotificationInterval: sql.NullString{
			String: event.NotificationInterval.String(),
			Valid:  event.NotificationInterval != 0,
		},
	}
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
