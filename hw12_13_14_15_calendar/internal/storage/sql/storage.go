package sqlstorage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewPostgresDB(ctx context.Context, dbCfg DBConf) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DBName,
		dbCfg.SSLMode,
	)

	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing config pgx: %s", err.Error())
	}

	conf.MaxConns = dbCfg.MaxConns
	conf.MinConns = dbCfg.MinConns
	conf.MaxConnIdleTime = dbCfg.MaxConnIdleTime
	conf.MaxConnLifetime = dbCfg.MaxConnLifetime

	db, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("error connecting pgx db: %w", err)
	}

	return db, nil
}
