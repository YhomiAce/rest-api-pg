package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/YhomiAce/rest-api-pg/cmd/api"
	"github.com/YhomiAce/rest-api-pg/config"
	"github.com/YhomiAce/rest-api-pg/db"
)

func main() {
	// Load configuration
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
	}
	// initialize database connection
	dbConn, err := db.NewPostgresSQL(
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	if err := initDB(dbConn); err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	server := api.NewAPIServer(":8080", dbConn, cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start API server: %v\n", err)
	}
}

func initDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}
	log.Println("Successfully connected to the database")
	return nil
}