package migrations

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic("failed to init driver: " + err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres",
		driver)
	if err != nil {
		panic("failed to init migrate instance: " + err.Error())
	}

	if err := m.Up(); !errors.Is(err, migrate.ErrNoChange) {
		panic("failed to migrate: " + err.Error())
	}
}
