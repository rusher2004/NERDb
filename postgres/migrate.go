package main

import (
	"log"
	"os"
	"path/filepath"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	dbURL, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		log.Fatal("POSTGRES_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create driver instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+filepath.Join("migrations"),
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}

	// if err := m.Force(1); err != nil {
	// 	log.Fatalf("failed to force migration: %v", err)
	// }

	if err := m.Up(); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
}
