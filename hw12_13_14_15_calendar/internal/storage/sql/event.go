package sqlstorage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type EventPostgres struct {
	db *pgxpool.Pool
}

func NewEventPostgres(db *pgxpool.Pool) *EventPostgres {
	return &EventPostgres{db: db}
}

func (e *EventPostgres) Create(ctx context.Context, event Event) (int, error) {
	panic("")
}

func (e *EventPostgres) Update(ctx context.Context, id int, event Event) (Event, error) {
	panic("")
}

func (e *EventPostgres) Delete(ctx context.Context, id int) (int, error) {
	panic("")
}
func (e *EventPostgres) GetAllByDay(ctx context.Context, data time.Time) ([]Event, error) {
	panic("")
}

func (e *EventPostgres) GetAllByWeek(ctx context.Context, data time.Time) ([]Event, error) {
	panic("")
}

func (e *EventPostgres) GetAllByMonth(ctx context.Context, data time.Time) ([]Event, error) {
	panic("")
}
