package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	"github.com/oriiyx/fritz/app/core/utils/logger"
)

var logFilePath = "logs/hard_reset.log"

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	l, err := logger.New(true, logFilePath)
	if err != nil {
		panic(err)
	}
	l.Info().Msg("Starting database hard reset process")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		l.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// Drop all tables, including the schema_migrations table
	dropQuery := `
	DO $$ DECLARE
		r RECORD;
	BEGIN
		-- Disable triggers
		SET session_replication_role = 'replica';
		
		-- Drop all tables in public schema
		FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
			EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
		END LOOP;

		-- Drop all sequences
		FOR r IN (SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema = 'public') LOOP
			EXECUTE 'DROP SEQUENCE IF EXISTS ' || quote_ident(r.sequence_name) || ' CASCADE';
		END LOOP;

		-- Drop all types
		FOR r IN (SELECT typname FROM pg_type 
				  JOIN pg_namespace ON pg_type.typnamespace = pg_namespace.oid 
				  WHERE nspname = 'public') LOOP
			EXECUTE 'DROP TYPE IF EXISTS ' || quote_ident(r.typname) || ' CASCADE';
		END LOOP;

		-- Enable triggers
		SET session_replication_role = 'origin';
	END $$;
	`

	// Execute the drop query
	_, err = db.ExecContext(ctx, dropQuery)
	if err != nil {
		l.Fatal().Stack().Err(err).Msg("Failed to drop database objects")
		return
	}

	l.Info().Msg("Successfully dropped all database objects")
	l.Info().Msg("You can now run 'make migrate-up' to recreate the database schema")
}
