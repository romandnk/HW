//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
	"time"
)

const url = "http://calendar:8080/api/v1/"

type CalendarSuite struct {
	suite.Suite
	ctx    context.Context
	dbConn *pgxpool.Pool
}

func (c *CalendarSuite) SetupSuite() {
	c.ctx = context.Background()

	connString := "postgres://test:1234@postgres:5432/calendar_db?sslmode=disable"

	conn, err := pgxpool.New(c.ctx, connString)
	c.Require().NoError(err)

	err = conn.Ping(c.ctx)
	c.Require().NoError(err)

	c.dbConn = conn
}

func (c *CalendarSuite) TearDownTest() {
	_, err := c.dbConn.Exec(c.ctx, "TRUNCATE events RESTART IDENTITY CASCADE")
	c.Require().NoError(err)
}

func (c *CalendarSuite) TestCreateEvent() {
	moscowTimeZone, err := time.LoadLocation("Europe/Moscow")
	c.Require().NoError(err)

	event := models.Event{
		Title:                "Harry Potter book",
		Date:                 time.Date(2023, 8, 20, 20, 30, 00, 00, moscowTimeZone), //nolint:gofumpt
		Duration:             time.Hour * 24,
		Description:          "new book for sale, no one has read it",
		UserID:               1,
		NotificationInterval: time.Hour,
		Scheduled:            false,
	}

	data := `{
		"title": "Harry Potter book",
		"date": "2023-08-20T20:30:00+03:00",
		"duration": "24h0m0s",
		"description": "new book for sale, no one has read it",
		"notification_interval": "1h0m0s",
		"user_id": 1
	}`

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, url+"events", bytes.NewReader([]byte(data)))
	c.Require().NoError(err)

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	c.Require().NoError(err)

	defer res.Body.Close()

	c.Require().Equal(http.StatusCreated, res.StatusCode)

	query := `
		SELECT * FROM events WHERE title = 'Harry Potter book'
	`

	var actualEvent models.Event
	err = c.dbConn.QueryRow(c.ctx, query).Scan(
		&actualEvent.ID,
		&actualEvent.Title,
		&actualEvent.Date,
		&actualEvent.Duration,
		&actualEvent.Description,
		&actualEvent.UserID,
		&actualEvent.NotificationInterval,
		&actualEvent.Scheduled,
	)
	c.Require().NoError(err)

	c.Require().True(actualEvent.ID != "")
	c.Require().Equal(event.Title, actualEvent.Title)
	c.Require().Equal(event.Date.Local(), actualEvent.Date)
	c.Require().Equal(event.Duration, actualEvent.Duration)
	c.Require().Equal(event.Description, actualEvent.Description)
	c.Require().Equal(event.UserID, actualEvent.UserID)
	c.Require().Equal(event.NotificationInterval, actualEvent.NotificationInterval)
	c.Require().Equal(event.Scheduled, actualEvent.Scheduled)
}

func (c *CalendarSuite) TestCreateEventEmptyTitle() {
	data := `{
		"title": "",
		"date": "2023-08-20T20:30:00+03:00",
		"duration": "24h0m0s",
		"description": "new book for sale, no one has read it",
		"notification_interval": "1h0m0s",
		"user_id": 1
	}`

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, url+"events", bytes.NewReader([]byte(data)))
	c.Require().NoError(err)

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	c.Require().NoError(err)

	defer res.Body.Close()

	c.Require().Equal(http.StatusInternalServerError, res.StatusCode)

	expectedAction := "create"
	expectedField := "title"
	expectedError := "title cannot be empty"
	expectedMessage := "error creating event"

	resBody, err := io.ReadAll(res.Body)
	c.Require().NoError(err)

	var response map[string]interface{}
	err = json.Unmarshal(resBody, &response)
	c.Require().NoError(err)

	action, ok := response["action"]
	c.Require().True(ok)
	c.Require().Equal(expectedAction, action)

	field, ok := response["field"]
	c.Require().True(ok)
	c.Require().Equal(field, expectedField)

	message, ok := response["message"]
	c.Require().True(ok)
	c.Require().Equal(message, expectedMessage)

	actualError, ok := response["error"]
	c.Require().True(ok)
	c.Require().Equal(actualError, expectedError)

	query := `
		SELECT * FROM events WHERE title = ''
	`

	var actualEvent models.Event
	err = c.dbConn.QueryRow(c.ctx, query).Scan(
		&actualEvent.ID,
		&actualEvent.Title,
		&actualEvent.Date,
		&actualEvent.Duration,
		&actualEvent.Description,
		&actualEvent.UserID,
		&actualEvent.NotificationInterval,
		&actualEvent.Scheduled,
	)
	c.Require().ErrorIs(err, pgx.ErrNoRows)
}

type eventsResponse struct {
	Total int            `json:"total"`
	Data  []eventDetails `json:"data"`
}

type eventDetails struct {
	ID                   string        `json:"id"`
	Title                string        `json:"title"`
	Date                 time.Time     `json:"date"`
	Duration             time.Duration `json:"duration"`
	Description          string        `json:"description"`
	UserID               int           `json:"user_id"`
	NotificationInterval time.Duration `json:"notification_interval"`
}

func (c *CalendarSuite) TestAllGetEventsByWeek() {
	startDate := time.Date(2023, 8, 1, 12, 0, 0, 0, time.UTC)

	events := []models.Event{
		{
			ID:                   uuid.New().String(),
			Title:                "title1",
			Date:                 startDate.AddDate(0, 0, 4),
			Duration:             time.Second,
			Description:          "description1",
			UserID:               1,
			NotificationInterval: time.Second,
			Scheduled:            false,
		},
		{
			ID:                   uuid.New().String(),
			Title:                "title2",
			Date:                 startDate.AddDate(0, 0, 6),
			Duration:             time.Second * 2,
			Description:          "description2",
			UserID:               2,
			NotificationInterval: time.Second * 2,
			Scheduled:            false,
		},
		{
			ID:                   uuid.New().String(),
			Title:                "title3",
			Date:                 startDate.AddDate(0, 0, 7),
			Duration:             time.Second * 3,
			Description:          "description3",
			UserID:               3,
			NotificationInterval: time.Second * 3,
			Scheduled:            false,
		},
	}

	queryInsertEvents := `
		INSERT INTO events (id, title, date, duration, description, user_id, notification_interval, scheduled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, event := range events {
		_, err := c.dbConn.Exec(c.ctx, queryInsertEvents,
			event.ID,
			event.Title,
			event.Date,
			event.Duration,
			event.Description,
			event.UserID,
			event.NotificationInterval,
			event.Scheduled)
		c.Require().NoError(err)
	}

	req, err := http.NewRequestWithContext(c.ctx,
		http.MethodGet,
		url+"events/week/2023-08-01T12:00:00Z",
		nil)
	c.Require().NoError(err)

	client := &http.Client{}
	res, err := client.Do(req)
	c.Require().NoError(err)

	defer res.Body.Close()

	c.Require().Equal(http.StatusOK, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	c.Require().NoError(err)

	var response eventsResponse
	err = json.Unmarshal(resBody, &response)
	c.Require().NoError(err)

	expectedTotal := 2
	expectedData := []eventDetails{
		{
			ID:                   events[0].ID,
			Title:                events[0].Title,
			Date:                 events[0].Date,
			Duration:             events[0].Duration,
			Description:          events[0].Description,
			UserID:               events[0].UserID,
			NotificationInterval: events[0].NotificationInterval,
		},
		{
			ID:                   events[1].ID,
			Title:                events[1].Title,
			Date:                 events[1].Date,
			Duration:             events[1].Duration,
			Description:          events[1].Description,
			UserID:               events[1].UserID,
			NotificationInterval: events[1].NotificationInterval,
		},
	}

	c.Require().Equal(expectedTotal, response.Total)
	c.Require().Equal(expectedData[0], response.Data[0])
	c.Require().Equal(expectedData[1], response.Data[1])
	c.Require().Len(response.Data, 2)
}

func (c *CalendarSuite) TestSender() {
	now := time.Now().UTC()

	events := []models.Event{
		{
			ID:                   uuid.New().String(),
			Title:                "title1",
			Date:                 now.Add(time.Second * 10),
			Duration:             time.Second,
			Description:          "description1",
			UserID:               1,
			NotificationInterval: 0,
			Scheduled:            false,
		},
		{
			ID:                   uuid.New().String(),
			Title:                "title2",
			Date:                 now.Add(time.Second * 8),
			Duration:             time.Second * 2,
			Description:          "description2",
			UserID:               2,
			NotificationInterval: time.Second * 2,
			Scheduled:            false,
		},
		{
			ID:                   uuid.New().String(),
			Title:                "title3",
			Date:                 now.Add(time.Second * 20),
			Duration:             time.Second * 3,
			Description:          "description3",
			UserID:               3,
			NotificationInterval: 0,
			Scheduled:            false,
		},
	}

	queryInsertEvents := `
		INSERT INTO events (id, title, date, duration, description, user_id, notification_interval, scheduled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, event := range events {
		_, err := c.dbConn.Exec(c.ctx, queryInsertEvents,
			event.ID,
			event.Title,
			event.Date,
			event.Duration,
			event.Description,
			event.UserID,
			event.NotificationInterval,
			event.Scheduled)
		c.Require().NoError(err)
	}

	time.Sleep(time.Second * 5)

	query := `
		SELECT scheduled
		FROM events 
		WHERE id = $1
	`

	scheduledEvents := make([]bool, 0)
	for _, event := range events {
		var scheduled bool

		err := c.dbConn.QueryRow(c.ctx, query, event.ID).Scan(&scheduled)
		c.Require().NoError(err)

		scheduledEvents = append(scheduledEvents, scheduled)
	}

	c.Require().True(scheduledEvents[0])
	c.Require().True(scheduledEvents[1])
	c.Require().False(scheduledEvents[2])
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, &CalendarSuite{})
}
