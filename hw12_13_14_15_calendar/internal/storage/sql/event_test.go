package sqlstorage

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
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

	columns := []string{"id"}
	rows := sqlmock.NewRows(columns).
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

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	rowsAfter := sqlmock.NewRows(columns).AddRow(eventAfter.ID, eventAfter.Title, eventAfter.Date, eventAfter.Duration,
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
	expectedError := fmt.Errorf("no event with id %s", eventBefore.ID)
	require.EqualError(t, err, expectedError.Error())
	require.Equal(t, models.Event{}, updatedEvent)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	id := uuid.New().String()

	ctx := context.Background()

	storage := NewStorageSQL(db)

	queryDelete := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, eventsTable)

	mock.ExpectExec(regexp.QuoteMeta(queryDelete)).WithArgs(driver.Value(id)).WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.Delete(ctx, id)
	require.NoError(t, err)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageDeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	id := uuid.New().String()

	ctx := context.Background()

	storage := NewStorageSQL(db)

	queryDelete := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, eventsTable)

	mock.ExpectExec(regexp.QuoteMeta(queryDelete)).WithArgs(driver.Value(id)).WillReturnResult(sqlmock.NewResult(1, 0))

	err = storage.Delete(ctx, id)
	expectedError := fmt.Errorf("no event with id %s", id)
	require.EqualError(t, err, expectedError.Error())

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageGetAllByDay(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

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

	storage := NewStorageSQL(db)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := sqlmock.NewRows(columns).
		AddRow("1", "Event 1", date, time.Hour, "Description 1", 1, time.Hour).
		AddRow("2", "Event 2", date, 2*time.Hour, "Description 2", 2, 2*time.Hour)

	queryGetByDay := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date = $1`, eventsTable)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetByDay)).WithArgs(driver.Value(date)).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByDay(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 2)
	require.Equal(t, expectedEvents, actualEvents)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageGetAllByDayEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	var expectedEvents []models.Event

	ctx := context.Background()

	storage := NewStorageSQL(db)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := sqlmock.NewRows(columns)

	queryGetByDay := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date = $1`, eventsTable)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetByDay)).WithArgs(driver.Value(date)).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByDay(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 0)
	require.Equal(t, expectedEvents, actualEvents)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageGetAllByWeek(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

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

	storage := NewStorageSQL(db)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := sqlmock.NewRows(columns).
		AddRow("1", "Event 1", date, time.Hour, "Description 1", 1, time.Hour).
		AddRow("2", "Event 2", date.AddDate(0, 0, 1), 2*time.Hour, "Description 2", 2, 2*time.Hour)

	queryGetByWeek := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date BETWEEN $1 AND $1 + INTERVAL '7 days'`, eventsTable)
	mock.ExpectQuery(regexp.QuoteMeta(queryGetByWeek)).WithArgs(driver.Value(date)).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByWeek(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 2)
	require.Equal(t, expectedEvents, actualEvents)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageGetAllByWeekEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	var expectedEvents []models.Event

	ctx := context.Background()

	storage := NewStorageSQL(db)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := sqlmock.NewRows(columns)

	queryGetByWeek := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date BETWEEN $1 AND $1 + INTERVAL '7 days'`, eventsTable)
	mock.ExpectQuery(regexp.QuoteMeta(queryGetByWeek)).WithArgs(driver.Value(date)).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByWeek(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 0)
	require.Equal(t, expectedEvents, actualEvents)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageGetAllByMonth(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

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

	storage := NewStorageSQL(db)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := sqlmock.NewRows(columns).
		AddRow("1", "Event 1", date, time.Hour, "Description 1", 1, time.Hour).
		AddRow("2", "Event 2", date.AddDate(0, 0, 1), 2*time.Hour, "Description 2", 2, 2*time.Hour)

	queryGetByMonth := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date BETWEEN $1 AND $1 + INTERVAL '1 month'`, eventsTable)
	mock.ExpectQuery(regexp.QuoteMeta(queryGetByMonth)).WithArgs(driver.Value(date)).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByMonth(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 2)
	require.Equal(t, expectedEvents, actualEvents)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}

func TestStorageGetAllByMonthEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	date := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)

	var expectedEvents []models.Event

	ctx := context.Background()

	storage := NewStorageSQL(db)

	columns := []string{"id", "title", "date", "duration", "description", "user_id", "notification_interval"}
	expectedRows := sqlmock.NewRows(columns)

	queryGetByMonth := fmt.Sprintf(`
		SELECT id, title, date, duration, description, user_id, notification_interval
		FROM %s 
		WHERE date BETWEEN $1 AND $1 + INTERVAL '1 month'`, eventsTable)
	mock.ExpectQuery(regexp.QuoteMeta(queryGetByMonth)).WithArgs(driver.Value(date)).WillReturnRows(expectedRows)

	actualEvents, err := storage.GetAllByMonth(ctx, date)
	require.NoError(t, err)
	require.Len(t, actualEvents, 0)
	require.Equal(t, expectedEvents, actualEvents)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal("there was unexpected result")
	}
}
