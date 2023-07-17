package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
)

func NewPostgresDB(ctx context.Context, cfg DBConf) (*sql.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("error loading pgx driver: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MinConns)
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return db, nil
}
