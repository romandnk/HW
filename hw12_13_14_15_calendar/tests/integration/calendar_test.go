//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/suite"
	"net"
	"net/http"
	"testing"
	"time"
)

const url = "http://calendar:8080"

type CalendarSuite struct {
	suite.Suite
	srv    *http.Server
	ctx    context.Context
	dbConn *pgxpool.Pool
}

func (c *CalendarSuite) SetupSuite() {
	srv := &http.Server{
		Addr: net.JoinHostPort("calendar", "8080"),
	}

	err := srv.ListenAndServe()
	c.Require().NoError(err)

	c.srv = srv
	c.ctx = context.Background()

	connString := "postgres://test:1234@postgres:5432/calendar_db?sslmode=disable"

	conn, err := pgxpool.New(c.ctx, connString)
	c.Require().NoError(err)

	err = conn.Ping(c.ctx)
	c.Require().NoError(err)

	c.dbConn = conn
}

func (c *CalendarSuite) TearDownSuite() {
	err := c.srv.Shutdown(c.ctx)
	c.Require().NoError(err)

	_, err = c.dbConn.Exec(c.ctx, "DROP TABLE events")
	c.Require().NoError(err)

	c.dbConn.Close()
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
		"Title": "Harry Potter book",
		"Date": "2023-08-20T20:30:00+03:00",
		"Duration": "24h0m0s",
		"Description": "new book for sale, no one has read it",
		"NotificationInterval": "1h0m0s",
	}`

	resp, err := http.Post(url+"/events", "application/json", bytes.NewReader([]byte(data)))
	c.Require().NoError(err)
	defer resp.Body.Close()

	c.Require().Equal(http.StatusCreated, resp.StatusCode)

	query := `
		SELECT * FROM events WHERE title = 'Harry Potter book'
	`
	rows, err := c.dbConn.Query(c.ctx, query)
	c.Require().NoError(err)

	var actualEvent models.Event
	err = rows.Scan(
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
	c.Require().Equal(event.Date.UTC(), actualEvent.Date)
	c.Require().Equal(event.Duration, actualEvent.Duration)
	c.Require().Equal(event.Description, actualEvent.Description)
	c.Require().Equal(event.UserID, actualEvent.UserID)
	c.Require().Equal(event.NotificationInterval, actualEvent.NotificationInterval)
	c.Require().Equal(event.Scheduled, actualEvent.Scheduled)
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, &CalendarSuite{})
}
