package base

import (
	"github.com/go-playground/validator/v10"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

// HandlerController provides common services for all handlers
type HandlerController struct {
	Logger    *zerolog.Logger
	Queries   *db.Queries
	Validator *validator.Validate
}

// HandlerControllerFactory creates HandlerController instances
type HandlerControllerFactory struct {
	baseLogger *zerolog.Logger
	queries    *db.Queries
	validator  *validator.Validate
}

// NewHandlerControllerFactory creates a factory from common services
func NewHandlerControllerFactory(logger *zerolog.Logger, queries *db.Queries, validator *validator.Validate) *HandlerControllerFactory {
	return &HandlerControllerFactory{
		baseLogger: logger,
		queries:    queries,
		validator:  validator,
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
		Logger:    logger,
		Queries:   f.queries,
		Validator: f.validator,
	}
}
