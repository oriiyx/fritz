package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	internalMiddleware "github.com/oriiyx/fritz/api/middleware"
	"github.com/oriiyx/fritz/api/middleware/requestlog"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/oriiyx/fritz/server/definitions"
	"github.com/oriiyx/fritz/utils/env"
	"github.com/rs/zerolog"
)

type Controller struct {
	Conf      *env.Conf
	Pool      *pgxpool.Pool
	Store     *sessions.CookieStore
	Router    *chi.Mux
	Logger    *zerolog.Logger
	Queries   *db.Queries
	Validator *validator.Validate
}

func (c *Controller) RegisterUses() {
	c.Router.Use(middleware.Logger)
	c.Router.Use(internalMiddleware.RequestID)
	c.Router.Use(internalMiddleware.JSONMiddleware)
	c.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},                      // Allow localhost and specific ports
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // Set to true if you need cookies for cross-site requests
		MaxAge:           300,  // Maximum value not to check for CORS again for a certain duration
	}))

	c.Logger.Info().Str("environment", c.Conf.Server.ENV).Msg("Environment")
}

func (c *Controller) RegisterRoutes() {
	c.Router.Route("/api/v1", func(r chi.Router) {
		definitionsHandler := definitions.NewDefinitionsHandler(c.Queries, c.Validator, c.Logger)
		r.Route("/definitions", func(definitions chi.Router) {
			definitions.Method(http.MethodGet, "/data-component-types", requestlog.NewHandler(definitionsHandler.GetDataComponentTypes, c.Logger))
		})
	})
}
