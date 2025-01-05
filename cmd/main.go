package main

import (
	"log"

	"github.com/sudarakas/edata/cmd/api"
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

	config := api.DefaultConfig()
	server := api.NewAPIServer(config, postgresDB)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
