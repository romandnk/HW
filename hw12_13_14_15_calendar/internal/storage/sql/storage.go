package sqlstorage

import (
	"context"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
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
	Create(ctx context.Context, event storage.Event) (string, error)
	Update(ctx context.Context, id string, event storage.Event) (storage.Event, error)
	Delete(ctx context.Context, id string) (string, error)
	GetAllByDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetAllByWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetAllByMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}
