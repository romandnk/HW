package sqlstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
)

func NewPostgresDB(ctx context.Context, cfg config.DBConfig) (PgxIface, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = int32(cfg.MaxConns)
	conf.MinConns = int32(cfg.MinConns)
	conf.MaxConnLifetime = cfg.MaxConnLifetime
	conf.MaxConnIdleTime = cfg.MaxConnIdleTime

	db, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
