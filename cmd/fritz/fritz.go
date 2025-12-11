package main

import (
	"context"
	"fmt"
	"log"

	"github.com/exaring/otelpgx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/oriiyx/fritz/app/core/kernel"
	"github.com/oriiyx/fritz/app/core/services"
	"github.com/oriiyx/fritz/app/core/utils/env"
	logger2 "github.com/oriiyx/fritz/app/core/utils/logger"
	"github.com/oriiyx/fritz/app/core/utils/validator"
	"github.com/rs/zerolog"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
const logFilePath = "var/logs/app.log"

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
	conf := env.New()

	l, err := logger2.New(conf.Server.Debug, logFilePath)
	if err != nil {
		log.Fatalf("Could not initiate logger with err: %v", err)
	}

	k := loadKernel(conf, l, ctx)
	if err = k.Start(ctx); err != nil {
		l.Fatal().Err(err).Msg("Failed to start the kernel.")
	}
}

func loadKernel(conf *env.Conf, l *zerolog.Logger, ctx context.Context) *kernel.Kernel {
	v := validator.New()

	l.Info().Str("env", conf.Server.ENV).Bool("debug", conf.Server.Debug).Msg("Starting application")
	traceLogger := logger2.NewTraceLogger(*l, conf.Server.Debug)
	m := logger2.MultiQueryTracer{
		Tracers: []pgx.QueryTracer{
			otelpgx.NewTracer(),
			traceLogger,
		},
	}

	l.Info().Msg("Initializing database connection")
	dbString := fmt.Sprintf(fmtDBString, conf.DB.Host, conf.DB.Username, conf.DB.Password, conf.DB.DBName, conf.DB.Port)

	l.Info().Str("host", conf.DB.Host).Int("port", conf.DB.Port).Str("database", conf.DB.DBName).Msg("Connecting to database")
	dbConfig, err := pgxpool.ParseConfig(dbString)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to parse database config")
	}
	dbConfig.ConnConfig.Tracer = &m
	pool, err := pgxpool.NewWithConfig(
		ctx,
		dbConfig,
	)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create database connection pool")
	}
	defer pool.Close()
	l.Info().Msg("Database connection established successfully")

	chiRouter := chi.NewRouter()

	k := kernel.New()

	if err = k.Registry().Register(services.Database, pool); err != nil {
		l.Fatal().Err(err).Msg("Failed to register database pool service.")
	}

	if err = k.Registry().Register(services.Logger, l); err != nil {
		l.Fatal().Err(err).Msg("Failed to register logger service.")
	}

	if err = k.Registry().Register(services.Validator, v); err != nil {
		l.Fatal().Err(err).Msg("Failed to register validator service.")
	}

	if err = k.Registry().Register(services.EnvConfig, conf); err != nil {
		l.Fatal().Err(err).Msg("Failed to register environment configurator service.")
	}

	if err = k.Registry().Register(services.Router, chiRouter); err != nil {
		l.Fatal().Err(err).Msg("Failed to register chi router service.")
	}

	return k
}
