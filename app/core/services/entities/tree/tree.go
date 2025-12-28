package tree

import (
	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/services/objects/definition_builder"
)

type Handler struct {
	*base.HandlerController
	entityBuilder *definition_builder.Builder
}

func New(ctrl *base.HandlerController) *Handler {
	eb := definition_builder.NewDefinitionsBuilder(ctrl.Logger, ctrl.DB, ctrl.CustomWriter)

	return &Handler{
		entityBuilder:     eb,
		HandlerController: ctrl,
	}
}
