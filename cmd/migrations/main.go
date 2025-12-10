package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/oriiyx/fritz/app/core/utils/logger"
)

var logFilePath = "var/logs/migration-app.log"

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	ctx := context.Background()
	l, err := logger.New(true, logFilePath)
	if err != nil {
		panic(err)
	}
	l.Info().Msg("Starting database migration process")
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to parse database config", err)
	}

	pool, err := pgxpool.NewWithConfig(
		ctx,
		dbConfig,
	)
	if err != nil {
		log.Fatal("Failed to create database connection pool", err)
	}
	defer pool.Close()
	l.Info().Msg("Database connection established successfully")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database ", err)
	}
	l.Info().Msg("Database connection opened successfully")

	// Create migration instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Point to your migration files. Here we're using local files, but it could be other sources.
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrations/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			l.Fatal().
				Err(err).
				Msg("Migration up failed")
		}
		l.Info().Msg("Migration up completed successfully")
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			l.Fatal().
				Err(err).
				Msg("Migration down failed")
		}
		l.Info().Msg("Migration down completed successfully")
	}
}
