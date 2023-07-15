package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"time"
)

func (s *Storage) Create(ctx context.Context, event models.Event) (string, error) {
	var id string

	query := fmt.Sprintf(`INSERT INTO %s (title, date, duration, description, user_id, notification_interval)
									VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, eventsTable)

	err := s.db.QueryRow(ctx, query,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationInterval,
	).Scan(&id)

	if err != nil {
		return id, fmt.Errorf("error creating event: %w", err)
	}

	return id, nil
}

func (s *Storage) Update(ctx context.Context, id string, event models.Event) (models.Event, error) {
	var updatedEvent models.Event

	query := fmt.Sprintf(`UPDATE %s SET title = $1, date = $2, duration = $3, description = $4, user_id = $5, notification_interval = $6 
          WHERE id = $7 RETURNING id, title, date, duration, description, user_id, notification_interval`, eventsTable)

	err := s.db.QueryRow(ctx, query,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationInterval,
		id,
	).Scan(&updatedEvent.ID, &updatedEvent.Title, &updatedEvent.Date, &updatedEvent.Duration,
		&updatedEvent.Description, &updatedEvent.UserID, &updatedEvent.NotificationInterval)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return updatedEvent, fmt.Errorf("no event with id %s: %w", id, err)
		}
		return updatedEvent, fmt.Errorf("error updating event: %w", err)
	}

	return updatedEvent, nil
}

func (s *Storage) Delete(ctx context.Context, id string) (string, error) {
	var deletedID string

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 RETURNING id`, eventsTable)

	err := s.db.QueryRow(ctx, query, id).Scan(&deletedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return deletedID, fmt.Errorf("no event with id %s: %w", id, err)
		}
		return deletedID, fmt.Errorf("error deleting event: %w", err)
	}

	return deletedID, nil
}

func (s *Storage) GetAllByDay(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event

	query := fmt.Sprintf(`SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s WHERE date = $1`, eventsTable)

	rows, err := s.db.Query(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("error selecting events by day: %w", err)
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
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) GetAllByWeek(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event

	query := fmt.Sprintf(`SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s WHERE date BETWEEN $1 AND $1 + INTERVAL '7 days'`, eventsTable)

	rows, err := s.db.Query(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("error selecting events by day: %w", err)
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
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) GetAllByMonth(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event

	query := fmt.Sprintf(`SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s WHERE date BETWEEN $1 AND $1 + INTERVAL '1 month'`, eventsTable)

	rows, err := s.db.Query(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("error selecting events by day: %w", err)
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
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}
