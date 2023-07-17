package sqlstorage

import (
	"database/sql"
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

type Storage struct {
	db *sql.DB
}

func NewStorageSQL(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}
