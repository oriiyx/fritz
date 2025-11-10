package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
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
	c.Logger.Info().Str("environment", c.Conf.Server.ENV).Msg("Environment")
}

func (c *Controller) RegisterRoutes() {

}
