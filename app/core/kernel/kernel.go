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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oriiyx/fritz/app/core/services"
	"github.com/oriiyx/fritz/app/core/utils/env"
	"github.com/rs/zerolog"
)

type Kernel struct {
	registry *Registry
	plugins  *Plugins
	hooks    *Hooks
}

func New() *Kernel {
	return &Kernel{
		registry: NewRegistry(),
		plugins:  NewPlugins(),
		hooks:    NewHooks(),
	}
}

// Registry returns the service registry
func (k *Kernel) Registry() *Registry {
	return k.registry
}

// Hooks returns the hook manager
func (k *Kernel) Hooks() *Hooks {
	return k.hooks
}

// RegisterPlugin adds a plugin to the kernel
func (k *Kernel) RegisterPlugin(plugin Plugin) {
	k.plugins.Register(plugin)
}

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

	conf := k.registry.MustGet(services.EnvConfig).(*env.Conf)
	l := k.registry.MustGet(services.Logger).(*zerolog.Logger)
	pool := k.registry.MustGet(services.Database).(*pgxpool.Pool)
	chiRouter := k.registry.MustGet(services.Router).(*chi.Mux)
	v := k.registry.MustGet(services.Validator).(*validator.Validate)

	// start the server
	routerController := Controller{
		Ctx:       ctx,
		Conf:      conf,
		Pool:      pool,
		Kernel:    k,
		Logger:    l,
		Router:    chiRouter,
		Validator: v,
	}

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

// Shutdown gracefully stops all plugins
func (k *Kernel) Shutdown(ctx context.Context) error {
	// Shutdown in reverse order
	for i := len(k.plugins.registered) - 1; i >= 0; i-- {
		plugin := k.plugins.registered[i]
		if err := plugin.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown plugin '%s': %w", plugin.Name(), err)
		}
	}

	return nil
}
