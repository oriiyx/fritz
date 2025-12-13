package base

import (
	"github.com/go-playground/validator/v10"
	"github.com/oriiyx/fritz/app/core/kernel"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

// HandlerController provides common services for all handlers
type HandlerController struct {
	Hooks     *kernel.Hooks
	Logger    *zerolog.Logger
	Queries   *db.Queries
	Validator *validator.Validate
}

// HandlerControllerFactory creates HandlerController instances
type HandlerControllerFactory struct {
	hooks      *kernel.Hooks
	queries    *db.Queries
	validator  *validator.Validate
	baseLogger *zerolog.Logger
}

// NewHandlerControllerFactory creates a factory from common services
func NewHandlerControllerFactory(logger *zerolog.Logger, queries *db.Queries, validator *validator.Validate, hooks *kernel.Hooks) *HandlerControllerFactory {
	return &HandlerControllerFactory{
		hooks:      hooks,
		queries:    queries,
		validator:  validator,
		baseLogger: logger,
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
		Hooks:     f.hooks,
		Logger:    logger,
		Queries:   f.queries,
		Validator: f.validator,
	}
}
