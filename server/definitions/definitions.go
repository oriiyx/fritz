package definitions

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/oriiyx/fritz/objects/definitions"
	"github.com/rs/zerolog"
)

type Handler struct {
	queries   *db.Queries
	validator *validator.Validate
	logger    *zerolog.Logger
}

func NewDefinitionsHandler(queries *db.Queries, validator *validator.Validate, logger *zerolog.Logger) *Handler {
	return &Handler{
		queries:   queries,
		validator: validator,
		logger:    logger,
	}
}

// GetDataComponentTypes returns all available data component types
func (h *Handler) GetDataComponentTypes(w http.ResponseWriter, r *http.Request) {
	components := definitions.GetAllDataComponents()
	_ = json.NewEncoder(w).Encode(components)
}
