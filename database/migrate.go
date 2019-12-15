package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"

	// Source driver, needed for some random reason
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// PerformMigrations migrates database to newest version using migration
// files under svenska-yle-bot/database/migrations
func (db *DB) PerformMigrations() error {
	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	//defer m.Close()
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}
