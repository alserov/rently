package postgres

import (
	"github.com/alserov/rently/carsharing/internal/db/migrations"
	"github.com/jmoiron/sqlx"
)

func MustConnect(dsn string) *sqlx.DB {
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to open db: " + err.Error())
	}

	if err = conn.Ping(); err != nil {
		panic("failed to ping db: " + err.Error())
	}

	migrations.Migrate(conn.DB)

	return conn
}
