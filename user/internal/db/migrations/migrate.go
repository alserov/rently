package migrations

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func Migrate(conn *sqlx.DB) {
	driver, err := mysql.WithInstance(conn.DB, &mysql.Config{})
	if err != nil {
		panic("failed to get driver: " + err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"mysql", driver)
	if err != nil {
		panic("failed to init migrate instance: " + err.Error())
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic("failed to migrate: " + err.Error())
	}
}
