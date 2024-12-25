package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresStorage(host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	// Format the DSN (Data Source Name) for PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Open a database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open connection:", err)
	}

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db, nil
}
