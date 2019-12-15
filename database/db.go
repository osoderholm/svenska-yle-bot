package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	// Postgres database driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB is a database instance
type DB struct {
	*sqlx.DB
}

// Open opens a Postgres database connection with given values
func Open(host, port, dbname, user, password string) (*DB, error) {
	db, err := sqlx.Connect("pgx",
		fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
			host, port, dbname, user, password))
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}
