package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YhomiAce/rest-api-pg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
	}
    m, err := migrate.New(
        "file://cmd/migrate/migrations",
        fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName,
		),
	)
    if err != nil {
		log.Fatalf("Failed to create migrate instance: %v\n", err)
	}
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v\n", err)
		}
		log.Println("Migrations applied successfully")
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v\n", err)
		}
		log.Println("Migrations rolled back successfully")
	}
}