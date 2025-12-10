package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/exaring/otelpgx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/oriiyx/fritz/app/core/api"
	"github.com/oriiyx/fritz/app/core/utils/env"
	logger2 "github.com/oriiyx/fritz/app/core/utils/logger"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
const logFilePath = "var/logs/app.log"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
	conf := env.New()

	l, err := logger2.New(conf.Server.Debug, logFilePath)
	if err != nil {
		log.Fatalf("Could not initiate logger with err: %v", err)
	}

	v := validator.New()
	ctx := context.Background()

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
		return
	}
	dbConfig.ConnConfig.Tracer = &m
	pool, err := pgxpool.NewWithConfig(
		ctx,
		dbConfig,
	)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create database connection pool")
		return
	}
	defer pool.Close()
	l.Info().Msg("Database connection established successfully")

	chiRouter := chi.NewRouter()
	routerController := api.Controller{
		Ctx:       ctx,
		Logger:    l,
		Validator: v,
		Pool:      pool,
		Router:    chiRouter,
		Conf:      conf,
	}

	routerController.RegisterUses()
	routerController.RegisterRoutes()

	addr := fmt.Sprintf("0.0.0.0:%d", conf.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      routerController.Router,
		ReadTimeout:  conf.Server.TimeoutRead,
		WriteTimeout: conf.Server.TimeoutWrite,
		IdleTimeout:  conf.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		l.Info().Msgf("Shutting down server %v", server.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), conf.Server.TimeoutIdle)
		defer cancel()

		if err = server.Shutdown(ctx); err != nil {
			l.Error().Err(err).Msg("Server shutdown failure")
		}

		if err == nil {
			pool.Close()
			l.Info().Msg("DB connection closed")
		}

		close(closed)
	}()

	l.Info().Str("address", addr).Int("port", conf.Server.Port).Msg("Starting HTTP server")
	if serverCloseErr := server.ListenAndServe(); serverCloseErr != nil && !errors.Is(serverCloseErr, http.ErrServerClosed) {
		l.Fatal().Err(serverCloseErr).Msg("Server startup failure")
	}

	<-closed
	l.Info().Msgf("Server shutdown successfully")
}
