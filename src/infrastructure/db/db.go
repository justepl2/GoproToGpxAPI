package db

import (
	"database/sql"
	"fmt"
	"os"
)

type Database struct {
	Conn *sql.DB
}

type Repository struct {
	db *Database
}

func NewDatabase() (*Database, error) {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")

	connStr := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPassword + " host=" + dbHost + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &Database{Conn: db}, nil
}

func (db *Database) NewRepository() *Repository {
	return &Repository{db: db}
}
