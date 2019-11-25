package pg

import (
	"database/sql"
	_ "github.com/lib/pq"
	"sync"
)

type PGConfig struct {
	DSN string `toml:"dsn"`
}

var (
	PGStore *sql.DB
	once    = new(sync.Once)
)

func InitPGStore(config PGConfig) (*sql.DB, error) {
	var err error
	once.Do(func() {
		PGStore, err = sql.Open("postgres", config.DSN)
	})

	return PGStore, err
}
