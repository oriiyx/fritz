package base

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oriiyx/fritz/app/core/kernel"
	"github.com/oriiyx/fritz/app/core/utils/env"
	"github.com/oriiyx/fritz/app/core/utils/rw"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

// HandlerController provides common services for all handlers
type HandlerController struct {
	DB           *pgxpool.Pool
	Conf         *env.Conf
	Hooks        *kernel.Hooks
	Logger       *zerolog.Logger
	Queries      *db.Queries
	Validator    *validator.Validate
	CustomWriter *rw.CustomWriter
}

// HandlerControllerFactory creates HandlerController instances
type HandlerControllerFactory struct {
	db           *pgxpool.Pool
	conf         *env.Conf
	hooks        *kernel.Hooks
	queries      *db.Queries
	validator    *validator.Validate
	baseLogger   *zerolog.Logger
	customWriter *rw.CustomWriter
}

// NewHandlerControllerFactory creates a factory from common services
func NewHandlerControllerFactory(
	logger *zerolog.Logger, queries *db.Queries,
	validator *validator.Validate, hooks *kernel.Hooks,
	db *pgxpool.Pool, cw *rw.CustomWriter, cfg *env.Conf,
) *HandlerControllerFactory {
	return &HandlerControllerFactory{
		db:           db,
		conf:         cfg,
		hooks:        hooks,
		queries:      queries,
		validator:    validator,
		baseLogger:   logger,
		customWriter: cw,
	}
}

// Create builds a HandlerController with an optional service name for scoped logging
func (f *HandlerControllerFactory) Create(serviceName string) *HandlerController {
	logger := f.baseLogger
	if serviceName != "" {
		scopedLogger := f.baseLogger.With().Str("service", serviceName).Logger()
		logger = &scopedLogger
	}

	return &HandlerController{
		DB:           f.db,
		Conf:         f.conf,
		Hooks:        f.hooks,
		Logger:       logger,
		Queries:      f.queries,
		Validator:    f.validator,
		CustomWriter: f.customWriter,
	}
}
