package sqlstorage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	MaxConns        int           // max connections in the pool
	MinConns        int           // min connections in the pool which must be opened
	MaxConnLifetime time.Duration // time after which db conn will be removed from the pool if there was no active use.
	MaxConnIdleTime time.Duration // time after which an inactive connection in the pool will be closed and deleted.
}

type PgxIface interface {
	Close()
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Storage struct {
	db PgxIface
}

func NewStorageSQL(db PgxIface) *Storage {
	return &Storage{
		db: db,
	}
}
