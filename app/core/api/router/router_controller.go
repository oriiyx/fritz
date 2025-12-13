package router

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oriiyx/fritz/app/core/kernel"
	"github.com/oriiyx/fritz/app/core/utils/env"
	"github.com/oriiyx/fritz/app/core/utils/writer"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

type Controller struct {
	DB           *pgxpool.Pool
	Ctx          context.Context
	Conf         *env.Conf
	Store        *sessions.CookieStore
	Kernel       *kernel.Kernel
	Router       *chi.Mux
	Logger       *zerolog.Logger
	Queries      *db.Queries
	Validator    *validator.Validate
	CustomWriter *writer.CustomWriter
}

func NewController(
	ctx context.Context, conf *env.Conf, db *pgxpool.Pool,
	store *sessions.CookieStore, kernel *kernel.Kernel, router *chi.Mux,
	l *zerolog.Logger, q *db.Queries, v *validator.Validate, cw *writer.CustomWriter,
) *Controller {
	return &Controller{
		DB:           db,
		Ctx:          ctx,
		Conf:         conf,
		Store:        store,
		Kernel:       kernel,
		Router:       router,
		Logger:       l,
		Queries:      q,
		Validator:    v,
		CustomWriter: cw,
	}
}
