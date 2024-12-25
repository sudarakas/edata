package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sudarakas/edata/config"
	"github.com/sudarakas/edata/db"
)

func main() {
	postgresDB, err := db.NewPostgresStorage(
		config.Envs.DBHost,
		config.Envs.DBPort,
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBName,
		config.Envs.DBSecure,
	)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to PostgreSQL database")
	}

	driver, err := postgres.WithInstance(postgresDB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args))-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up applied successfully")
	} else if cmd == "down" {
		if err = m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down rolled back successfully")
	} else {
		log.Fatal("Invalid command")
	}

}
