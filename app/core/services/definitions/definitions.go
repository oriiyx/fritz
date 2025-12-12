package definitions

import (
	"encoding/json"
	"net/http"

	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
)

type Handler struct {
	*base.HandlerController
}

func NewDefinitionsHandler(ctrl *base.HandlerController) *Handler {
	return &Handler{
		HandlerController: ctrl,
	}
}

// GetDataComponentTypes returns all available data component types
func (h *Handler) GetDataComponentTypes(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(definitions.GetAllDataComponents())
}
