package sqlstorage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestStorageCreateEvent(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	event := models.Event{
		ID:                   uuid.New().String(),
		Title:                "test title",
		Date:                 time.Now(),
		Duration:             time.Second,
		Description:          "test description",
		UserID:               4,
		NotificationInterval: time.Second,
	}

	columns := []string{"id"}
	rows := pgxmock.NewRows(columns).AddRow(event.ID)

	ctx := context.Background()

	query := fmt.Sprintf(`
		INSERT INTO %s (id, title, date, duration, description, user_id, notification_interval)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`, eventsTable)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(
		event.ID,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationInterval).WillReturnRows(rows)

	storage := NewStorageSQL(mock)

	id, err := storage.CreateEvent(ctx, event)

	require.NoError(t, err)
	require.Equal(t, event.ID, id)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err, "there was unexpected result")
}

func TestStorageUpdateEvent(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

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

	storage := NewStorageSQL(mock)

	columnsSelect := []string{"title", "date", "duration", "description", "user_id", "notification_interval"}
	rows := pgxmock.NewRows(columnsSelect).AddRow(eventBefore.Title, eventBefore.Date, eventBefore.Duration,
		eventBefore.Description, eventBefore.UserID, eventBefore.NotificationInterval)

	querySelect := fmt.Sprintf(`
			SELECT 
			    title, 
			    date, 
			    duration,
			    description,
			    user_id,
			    notification_interval
			FROM %s 
			WHERE id = $1`, eventsTable)

	columnsUpdate := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	rowsAfter := pgxmock.NewRows(columnsUpdate).AddRow(eventAfter.ID, eventAfter.Title, eventAfter.Date,
		eventAfter.Duration, eventAfter.Description, eventAfter.UserID, eventAfter.NotificationInterval)

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

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelect)).WithArgs(eventAfter.ID).WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).WithArgs(
		eventAfter.Title,
		eventAfter.Date,
		eventAfter.Duration,
		eventAfter.Description,
		eventAfter.UserID,
		eventAfter.NotificationInterval,
		eventAfter.ID).WillReturnRows(rowsAfter)
	mock.ExpectCommit()

	updatedEvent, err := storage.UpdateEvent(ctx, eventBefore.ID, eventAfter)
	require.NoError(t, err)
	require.Equal(t, eventAfter, updatedEvent)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageUpdateEventError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

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

	storage := NewStorageSQL(mock)

	querySelect := fmt.Sprintf(`
			SELECT 
			    title, 
			    date, 
			    duration,
			    description,
			    user_id,
			    notification_interval
			FROM %s 
			WHERE id = $1`, eventsTable)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelect)).WithArgs(eventAfter.ID).WillReturnError(pgx.ErrNoRows)
	mock.ExpectRollback()

	updatedEvent, err := storage.UpdateEvent(ctx, eventBefore.ID, eventAfter)
	expectedError := fmt.Errorf("no event with id %s", eventBefore.ID)
	require.EqualError(t, err, expectedError.Error())
	require.Equal(t, models.Event{}, updatedEvent)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageDeleteEvent(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := uuid.New().String()

	ctx := context.Background()

	storage := NewStorageSQL(mock)

	queryDelete := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, eventsTable)

	mock.ExpectExec(regexp.QuoteMeta(queryDelete)).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = storage.DeleteEvent(ctx, id)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageDeleteEventError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := uuid.New().String()

	ctx := context.Background()

	storage := NewStorageSQL(mock)

	queryDelete := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, eventsTable)

	mock.ExpectExec(regexp.QuoteMeta(queryDelete)).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = storage.DeleteEvent(ctx, id)
	expectedError := fmt.Errorf("no event with id %s", id)
	require.EqualError(t, err, expectedError.Error())

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageGetAllByDayEvents(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	expectedEvents := []models.Event{
		{
			ID:                   "1",
			Title:                "Event 1",
			Date:                 date,
			Duration:             time.Hour,
			Description:          "Description 1",
			UserID:               1,
			NotificationInterval: time.Hour,
		},
		{
			ID:                   "2",
			Title:                "Event 2",
			Date:                 date,
			Duration:             2 * time.Hour,
			Description:          "Description 2",
			UserID:               2,
			NotificationInterval: 2 * time.Hour,
		},
	}

	ctx := context.Background()

	storage := NewStorageSQL(mock)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := pgxmock.NewRows(columns).
		AddRow("1", "Event 1", date, time.Hour, "Description 1", 1, time.Hour).
		AddRow("2", "Event 2", date, 2*time.Hour, "Description 2", 2, 2*time.Hour)

	queryGetByDay := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s
		WHERE date = $1`, eventsTable)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetByDay)).WithArgs(date).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByDayEvents(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 2)
	require.ElementsMatch(t, expectedEvents, actualEvents)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageGetAllByDayEventsEmpty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	var expectedEvents []models.Event

	ctx := context.Background()

	storage := NewStorageSQL(mock)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := pgxmock.NewRows(columns)

	queryGetByDay := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s
		WHERE date = $1`, eventsTable)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetByDay)).WithArgs(date).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByDayEvents(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 0)
	require.ElementsMatch(t, expectedEvents, actualEvents)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageGetAllByWeekEvents(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	expectedEvents := []models.Event{
		{
			ID:                   "1",
			Title:                "Event 1",
			Date:                 date,
			Duration:             time.Hour,
			Description:          "Description 1",
			UserID:               1,
			NotificationInterval: time.Hour,
		},
		{
			ID:                   "2",
			Title:                "Event 2",
			Date:                 date.AddDate(0, 0, 1),
			Duration:             2 * time.Hour,
			Description:          "Description 2",
			UserID:               2,
			NotificationInterval: 2 * time.Hour,
		},
	}

	ctx := context.Background()

	storage := NewStorageSQL(mock)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := pgxmock.NewRows(columns).
		AddRow("1", "Event 1", date, time.Hour, "Description 1", 1, time.Hour).
		AddRow("2", "Event 2", date.AddDate(0, 0, 1), 2*time.Hour, "Description 2", 2, 2*time.Hour)

	queryGetByWeek := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s
		WHERE date BETWEEN $1 AND $1 + INTERVAL '7 days'`, eventsTable)
	mock.ExpectQuery(regexp.QuoteMeta(queryGetByWeek)).WithArgs(date).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByWeekEvents(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 2)
	require.ElementsMatch(t, expectedEvents, actualEvents)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageGetAllByWeekEventsEmpty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	var expectedEvents []models.Event

	ctx := context.Background()

	storage := NewStorageSQL(mock)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := pgxmock.NewRows(columns)

	queryGetByWeek := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s
		WHERE date BETWEEN $1 AND $1 + INTERVAL '7 days'`, eventsTable)
	mock.ExpectQuery(regexp.QuoteMeta(queryGetByWeek)).WithArgs(date).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByWeekEvents(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 0)
	require.ElementsMatch(t, expectedEvents, actualEvents)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageGetAllByMonthEvents(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	expectedEvents := []models.Event{
		{
			ID:                   "1",
			Title:                "Event 1",
			Date:                 date,
			Duration:             time.Hour,
			Description:          "Description 1",
			UserID:               1,
			NotificationInterval: time.Hour,
		},
		{
			ID:                   "2",
			Title:                "Event 2",
			Date:                 date.AddDate(0, 0, 1),
			Duration:             2 * time.Hour,
			Description:          "Description 2",
			UserID:               2,
			NotificationInterval: 2 * time.Hour,
		},
	}

	ctx := context.Background()

	storage := NewStorageSQL(mock)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := pgxmock.NewRows(columns).
		AddRow("1", "Event 1", date, time.Hour, "Description 1", 1, time.Hour).
		AddRow("2", "Event 2", date.AddDate(0, 0, 1), 2*time.Hour, "Description 2", 2, 2*time.Hour)

	queryGetByMonth := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s
		WHERE date BETWEEN $1 AND $1 + INTERVAL '1 month'`, eventsTable)
	mock.ExpectQuery(regexp.QuoteMeta(queryGetByMonth)).WithArgs(date).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByMonthEvents(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 2)
	require.ElementsMatch(t, expectedEvents, actualEvents)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestStorageGetAllByMonthEventsEmpty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	var expectedEvents []models.Event

	ctx := context.Background()

	storage := NewStorageSQL(mock)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := pgxmock.NewRows(columns)

	queryGetByMonth := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s
		WHERE date BETWEEN $1 AND $1 + INTERVAL '1 month'`, eventsTable)
	mock.ExpectQuery(regexp.QuoteMeta(queryGetByMonth)).WithArgs(date).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByMonthEvents(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 0)
	require.ElementsMatch(t, expectedEvents, actualEvents)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
