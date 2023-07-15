package sqlstorage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
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

//go:generate mockgen -source=storage.go -destination=.mock/mock.go sqlstorage
type StoreEvent interface {
	Create(ctx context.Context, event models.Event) (string, error)
	Update(ctx context.Context, id string, event models.Event) (models.Event, error)
	Delete(ctx context.Context, id string) (string, error)
	GetAllByDay(ctx context.Context, date time.Time) ([]models.Event, error)
	GetAllByWeek(ctx context.Context, date time.Time) ([]models.Event, error)
	GetAllByMonth(ctx context.Context, date time.Time) ([]models.Event, error)
}

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		db: db,
	}
}
