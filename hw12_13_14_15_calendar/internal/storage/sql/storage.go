package sqlstorage

import (
	"context"
	"time"
)

const eventsTable = "events"

type DBConf struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	MaxConns        int32         // max connections in the pool
	MinConns        int32         // min connections in the pool which must be opened
	MaxConnLifetime time.Duration // time after which db conn will be closed and removed from the pool if there was no active use.
	MaxConnIdleTime time.Duration // time after which an inactive connection in the pool will be closed and deleted.
}

type Event interface {
	Create(ctx context.Context, event Event) (int, error)
	Update(ctx context.Context, id int, event Event) (Event, error)
	Delete(ctx context.Context, id int) (int, error)
	GetAllByDay(ctx context.Context, data time.Time) ([]Event, error)
	GetAllByWeek(ctx context.Context, data time.Time) ([]Event, error)
	GetAllByMonth(ctx context.Context, data time.Time) ([]Event, error)
}
