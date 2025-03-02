package main

import (
	"database/sql"
	"fmt"
	"log"

	"blazperic/radionica/config"
	"blazperic/radionica/internal/utils"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	migrationsDir := "./migrations"
	if err := utils.RunMigrations(db, migrationsDir); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
