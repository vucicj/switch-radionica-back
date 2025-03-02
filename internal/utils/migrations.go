package utils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type Migration struct {
	Name    string
	Content string
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Ensure migrations table exists (already created in 0001)
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
			}
			migrations = append(migrations, Migration{
				Name:    file.Name(),
				Content: string(content),
			})
		}
	}

	// Sort migrations by name (e.g., 0001_..., 0002_...) to ensure order
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if anything fails

	for _, migration := range migrations {
		// Check if migration was already applied
		var exists bool
		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE name = $1)", migration.Name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %v", migration.Name, err)
		}
		if exists {
			log.Printf("Skipping already applied migration: %s", migration.Name)
			continue
		}

		// Apply migration
		_, err = tx.Exec(migration.Content)
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", migration.Name, err)
		}

		// Record migration in tracking table
		_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", migration.Name)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %v", migration.Name, err)
		}
		log.Printf("Applied migration: %s", migration.Name)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migrations: %v", err)
	}
	log.Println("All migrations applied successfully")
	return nil
}
