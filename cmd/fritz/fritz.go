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
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/oriiyx/fritz/app/core/api/router"
	"github.com/oriiyx/fritz/app/core/kernel"
	"github.com/oriiyx/fritz/app/core/services"
	"github.com/oriiyx/fritz/app/core/services/entities/adapters"
	"github.com/oriiyx/fritz/app/core/utils/env"
	logger2 "github.com/oriiyx/fritz/app/core/utils/logger"
	"github.com/oriiyx/fritz/app/core/utils/rw"
	internalValidator "github.com/oriiyx/fritz/app/core/utils/validator"
	db "github.com/oriiyx/fritz/database/generated"
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

	// Start kernel (initializes plugins)
	if err = k.Start(ctx); err != nil {
		l.Fatal().Err(err).Msg("Failed to start the kernel.")
	}

	// Now set up application (after kernel is initialized)
	startApplication(ctx, k, conf, l)
}

func loadKernel(conf *env.Conf, l *zerolog.Logger, ctx context.Context) *kernel.Kernel {
	v := internalValidator.New()

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
	// defer pool.Close()
	l.Info().Msg("Database connection established successfully")

	chiRouter := chi.NewRouter()
	store := createCookieStore(conf)
	queries := db.New(pool)
	customWriter := rw.New(l)

	adapters.LoadAll(queries)

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

	if err = k.Registry().Register(services.CookieStore, store); err != nil {
		l.Fatal().Err(err).Msg("Failed to register cookie store service.")
	}

	if err = k.Registry().Register(services.Queries, queries); err != nil {
		l.Fatal().Err(err).Msg("Failed to register queries service.")
	}

	if err = k.Registry().Register(services.CustomWriter, customWriter); err != nil {
		l.Fatal().Err(err).Msg("Failed to register custom writer service.")
	}

	return k
}

func startApplication(ctx context.Context, k *kernel.Kernel, conf *env.Conf, l *zerolog.Logger) {
	// Fetch all required services
	chiRouter := k.Registry().MustGet(services.Router).(*chi.Mux)
	queries := k.Registry().MustGet(services.Queries).(*db.Queries)
	pool := k.Registry().MustGet(services.Database).(*pgxpool.Pool)
	v := k.Registry().MustGet(services.Validator).(*validator.Validate)
	store := k.Registry().MustGet(services.CookieStore).(*sessions.CookieStore)
	cw := k.Registry().MustGet(services.CustomWriter).(*rw.CustomWriter)

	// Create router controller
	routerController := router.NewController(ctx, conf, pool, store, k, chiRouter, l, queries, v, cw)
	routerController.RegisterUses()
	routerController.RegisterRoutes()

	if err := k.Registry().Register(services.Controller, routerController); err != nil {
		l.Fatal().Err(err).Msg("Failed to register router controller service.")
	}

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

		l.Info().Msg("Shutting down Kernel")
		err := k.Shutdown(ctx)
		if err != nil {
			l.Error().Stack().Err(err).Msg("Failed to shutdown Kernel")
		}

		l.Info().Msgf("Shutting down Server %v", server.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), conf.Server.TimeoutIdle)
		defer cancel()

		err = server.Shutdown(ctx)
		if err != nil {
			l.Error().Err(err).Msg("Server shutdown failure")
		}

		if err == nil {
			pool.Close()
			l.Info().Msg("Database connection closed")
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

func createCookieStore(conf *env.Conf) *sessions.CookieStore {
	store := sessions.NewCookieStore(conf.Server.Secret)
	store.MaxAge(conf.Auth.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode

	return store
}
