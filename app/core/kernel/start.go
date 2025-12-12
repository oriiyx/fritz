package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oriiyx/fritz/app/core/services"
	"github.com/oriiyx/fritz/app/core/utils/env"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

// Start initializes kernel and all plugins
func (k *Kernel) Start(ctx context.Context) error {
	// Sort plugins by priority
	sort.Slice(k.plugins.registered, func(i, j int) bool {
		return k.plugins.registered[i].Priority() < k.plugins.registered[j].Priority()
	})

	// Initialize each plugin
	for _, plugin := range k.plugins.registered {
		if err := plugin.Initialize(ctx, k); err != nil {
			return fmt.Errorf("failed to initialize plugin '%s': %w", plugin.Name(), err)
		}

		k.plugins.initialized[plugin.Name()] = plugin
	}

	// Fetch all required services
	l := k.registry.MustGet(services.Logger).(*zerolog.Logger)
	conf := k.registry.MustGet(services.EnvConfig).(*env.Conf)
	chiRouter := k.registry.MustGet(services.Router).(*chi.Mux)
	queries := k.registry.MustGet(services.Queries).(*db.Queries)
	pool := k.registry.MustGet(services.Database).(*pgxpool.Pool)
	v := k.registry.MustGet(services.Validator).(*validator.Validate)
	store := k.registry.MustGet(services.CookieStore).(*sessions.CookieStore)

	// Create
	routerController := NewController(ctx, conf, pool, store, k, chiRouter, l, queries, v)
	routerController.RegisterUses()
	routerController.RegisterRoutes()

	if err := k.registry.Register(services.Controller, routerController); err != nil {
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
	return nil
}
