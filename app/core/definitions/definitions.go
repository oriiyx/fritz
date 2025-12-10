package definitions

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/oriiyx/fritz/app/core/objects/definitions"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

type Handler struct {
	queries   *db.Queries
	validator *validator.Validate
	logger    *zerolog.Logger
}

func NewDefinitionsHandler(queries *db.Queries, validator *validator.Validate, logger *zerolog.Logger) *Handler {
	loggerWithService := logger.With().Str("service", "definitions").Logger()

	return &Handler{
		queries:   queries,
		validator: validator,
		logger:    &loggerWithService,
	}
}

// GetDataComponentTypes returns all available data component types
func (h *Handler) GetDataComponentTypes(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(definitions.GetAllDataComponents())
}
