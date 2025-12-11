package kernel

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oriiyx/fritz/app/core/utils/env"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

type Controller struct {
	Ctx       context.Context
	Conf      *env.Conf
	Pool      *pgxpool.Pool
	Store     *sessions.CookieStore
	Kernel    *Kernel
	Router    *chi.Mux
	Logger    *zerolog.Logger
	Queries   *db.Queries
	Validator *validator.Validate
}
