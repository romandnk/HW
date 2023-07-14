package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type EventPostgres struct {
	db *pgxpool.Pool
}

func NewEventPostgres(db *pgxpool.Pool) *EventPostgres {
	return &EventPostgres{db: db}
}

func (e *EventPostgres) Create(ctx context.Context, event storage.Event) (string, error) {
	var id string

	query := fmt.Sprintf(`INSERT INTO %s (title, date, duration, description, user_id, notification_interval)
									VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, eventsTable)

	err := e.db.QueryRow(ctx, query,
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

func (e *EventPostgres) Update(ctx context.Context, id string, event storage.Event) (storage.Event, error) {
	var updatedEvent storage.Event

	query := fmt.Sprintf(`UPDATE %s SET title = $1, date = $2, duration = $3, description = $4, user_id = $5, notification_interval = $6 
          WHERE id = $7 RETURNING id, title, date, duration, description, user_id, notification_interval`, eventsTable)

	err := e.db.QueryRow(ctx, query,
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
			return updatedEvent, fmt.Errorf("no event with id %d: %w", id, err)
		}
		return updatedEvent, fmt.Errorf("error updating event: %w", err)
	}

	return updatedEvent, nil
}

func (e *EventPostgres) Delete(ctx context.Context, id string) (string, error) {
	var deletedID string

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 RETURNING id`, eventsTable)

	err := e.db.QueryRow(ctx, query, id).Scan(&deletedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return deletedID, fmt.Errorf("no event with id %s: %w", id, err)
		}
		return deletedID, fmt.Errorf("error deleting event: %w", err)
	}

	return deletedID, nil
}

func (e *EventPostgres) GetAllByDay(ctx context.Context, data time.Time) ([]storage.Event, error) {
	panic("")
}

func (e *EventPostgres) GetAllByWeek(ctx context.Context, data time.Time) ([]storage.Event, error) {
	panic("")
}

func (e *EventPostgres) GetAllByMonth(ctx context.Context, data time.Time) ([]storage.Event, error) {
	panic("")
}
