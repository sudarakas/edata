package main

import (
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sudarakas/edata/config"
	"github.com/sudarakas/edata/db"
)

func main() {
	// Initialize the PostgreSQL database connection
	postgresDB, err := db.NewPostgresStorage(
		config.Envs.DBHost,
		config.Envs.DBPort,
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBName,
		config.Envs.DBSecure,
	)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL database:", err)
	}
	log.Println("Connected to PostgreSQL database")

	// Create a migration driver for PostgreSQL
	driver, err := postgres.WithInstance(postgresDB, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create PostgreSQL migration driver:", err)
	}

	// Initialize the migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to initialize migration instance:", err)
	}

	// Process command-line arguments
	if len(os.Args) < 2 {
		log.Fatal("Missing command argument. Valid commands are: up, down, force")
	}
	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down rolled back successfully")

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Missing version number for the force command. Usage: force <version>")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version number for the force command: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("Force migration to version %d failed: %v", version, err)
		}
		log.Printf("Force migration to version %d applied successfully", version)

	default:
		log.Fatalf("Invalid command: %s. Valid commands are: up, down, force", cmd)
	}
}
