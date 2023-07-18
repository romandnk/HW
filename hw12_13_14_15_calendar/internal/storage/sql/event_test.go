package sqlstorage

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestStorageCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	event := models.Event{
		ID:                   uuid.New().String(),
		Title:                "test title",
		Date:                 time.Now(),
		Duration:             time.Second,
		Description:          "test description",
		UserID:               4,
		NotificationInterval: time.Second,
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(event.ID)

	ctx := context.Background()

	query := fmt.Sprintf(`
		INSERT INTO %s (id, title, date, duration, description, user_id, notification_interval)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`, eventsTable)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(
		driver.Value(event.ID),
		driver.Value(event.Title),
		driver.Value(event.Date),
		driver.Value(event.Duration),
		driver.Value(event.Description),
		driver.Value(event.UserID),
		driver.Value(event.NotificationInterval)).WillReturnRows(rows)

	storage := NewStorageSQL(db)

	id, err := storage.Create(ctx, event)

	require.NoError(t, err)
	require.Equal(t, event.ID, id)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	eventBefore := models.Event{
		ID:                   uuid.New().String(),
		Title:                "test title update before",
		Date:                 time.Now(),
		Duration:             time.Second,
		Description:          "test description update before",
		UserID:               4,
		NotificationInterval: time.Second,
	}

	eventAfter := models.Event{
		ID:                   eventBefore.ID,
		Title:                "test title update after",
		Date:                 eventBefore.Date,
		Duration:             time.Second,
		Description:          "test description update after",
		UserID:               4,
		NotificationInterval: time.Second,
	}

	ctx := context.Background()

	storage := NewStorageSQL(db)

	rowsAfter := sqlmock.NewRows([]string{
		"id", "title", "date", "duration", "description", "user_id", "notification_interval",
	}).AddRow(eventAfter.ID, eventAfter.Title, eventAfter.Date, eventAfter.Duration,
		eventAfter.Description, eventAfter.UserID, eventAfter.NotificationInterval)

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

	mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).WithArgs(
		driver.Value(eventAfter.Title),
		driver.Value(eventAfter.Date),
		driver.Value(eventAfter.Duration),
		driver.Value(eventAfter.Description),
		driver.Value(eventAfter.UserID),
		driver.Value(eventAfter.NotificationInterval),
		driver.Value(eventAfter.ID)).WillReturnRows(rowsAfter)

	updatedEvent, err := storage.Update(ctx, eventBefore.ID, eventAfter)
	require.NoError(t, err)
	require.Equal(t, eventAfter, updatedEvent)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageUpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	eventBefore := models.Event{
		ID:                   uuid.New().String(),
		Title:                "test title update before",
		Date:                 time.Now(),
		Duration:             time.Second,
		Description:          "test description update before",
		UserID:               4,
		NotificationInterval: time.Second,
	}

	eventAfter := models.Event{
		ID:                   eventBefore.ID,
		Title:                "test title update after",
		Date:                 eventBefore.Date,
		Duration:             time.Second,
		Description:          "test description update after",
		UserID:               4,
		NotificationInterval: time.Second,
	}

	ctx := context.Background()

	storage := NewStorageSQL(db)

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

	mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).WithArgs(
		driver.Value(eventAfter.Title),
		driver.Value(eventAfter.Date),
		driver.Value(eventAfter.Duration),
		driver.Value(eventAfter.Description),
		driver.Value(eventAfter.UserID),
		driver.Value(eventAfter.NotificationInterval),
		driver.Value(eventAfter.ID)).WillReturnError(pgx.ErrNoRows)

	updatedEvent, err := storage.Update(ctx, eventBefore.ID, eventAfter)
	expectedError := fmt.Errorf("no event with id %s: %w", eventBefore.ID, pgx.ErrNoRows)
	require.Equal(t, err, expectedError)
	require.Equal(t, models.Event{}, updatedEvent)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}
