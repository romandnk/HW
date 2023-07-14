package sqlstorage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg DBConf) (*pgxpool.Pool, error) {
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
		return nil, fmt.Errorf("error parsing config pgx: %s", err.Error())
	}

	conf.MaxConns = cfg.MaxConns
	conf.MinConns = cfg.MinConns
	conf.MaxConnIdleTime = cfg.MaxConnIdleTime
	conf.MaxConnLifetime = cfg.MaxConnLifetime

	db, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("error connecting pgx db: %w", err)
	}

	return db, nil
}
