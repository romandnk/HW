//go:build integration

package intergration_test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"net"
	"net/http"
	"os"
	"testing"
)

type CalendarSuite struct {
	suite.Suite
	srv    *http.Server
	ctx    context.Context
	dbConn *pgxpool.Pool
}

func (c *CalendarSuite) SetupSuite() {
	host := os.Getenv("TEST_HOST")
	port := os.Getenv("TEST_PORT")

	srv := &http.Server{
		Addr: net.JoinHostPort(host, port),
	}

	err := srv.ListenAndServe()
	c.Require().NoError(err)

	c.srv = srv
	c.ctx = context.Background()

	//connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	//	os.Getenv("TEST_POSTGRES_USERNAME"),
	//	os.Getenv("TEST_POSTGRES_PASSWORD"),
	//	os.Getenv("TEST_POSTGRES_HOST"),
	//	os.Getenv("TEST_POSTGRES_PORT"),
	//	os.Getenv("TEST_POSTGRES_DBNAME"),
	//)

	conn, err := pgxpool.New(c.ctx, connString)
	c.Require().NoError(err)

	err = conn.Ping(c.ctx)
	c.Require().NoError(err)

	_, err = conn.Exec(c.ctx, `
		CREATE TABLE test_events (
    		id VARCHAR(36) PRIMARY KEY,
    		title VARCHAR(255) NOT NULL,
    		date TIMESTAMPTZ NOT NULL,
    		duration INTERVAL HOUR TO SECOND NOT NULL,
    		description TEXT,
    		user_id INTEGER NOT NULL,
    		notification_interval INTERVAL,
    		scheduled boolean DEFAULT FALSE NOT NULL
		)`)
	c.Require().NoError(err)

	c.dbConn = conn
}

func (c *CalendarSuite) TearDownSuite() {
	err := c.srv.Shutdown(c.ctx)
	c.Require().NoError(err)

	_, err = c.dbConn.Exec(c.ctx, "DROP TABLE test_events")
	c.Require().NoError(err)

	c.dbConn.Close()
}

func (c *CalendarSuite) TearDownTest() {
	_, err := c.dbConn.Exec(c.ctx, "TRUNCATE test_events RESTART IDENTITY CASCADE")
	c.Require().NoError(err)
}

func (c *CalendarSuite) TestCreateEvent() {

}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, &CalendarSuite{})
}
