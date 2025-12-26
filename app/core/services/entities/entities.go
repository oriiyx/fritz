package entities

import (
	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/services/objects/definition_builder"
)

const (
	DefinitionIDKey = "definition_id"
)

type Handler struct {
	*base.HandlerController

	entityBuilder *definition_builder.Builder
}

func NewEntitiesHandler(ctrl *base.HandlerController) *Handler {
	eb := definition_builder.NewDefinitionsBuilder(ctrl.Logger, ctrl.DB, ctrl.CustomWriter)

	return &Handler{
		HandlerController: ctrl,
		entityBuilder:     eb,
	}
}
