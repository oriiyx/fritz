package entities

import (
	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/services/objects/entity_builder"
)

const (
	ClassIDKey = "class_id"
)

type Handler struct {
	*base.HandlerController

	entityBuilder *entity_builder.EntityBuilder
}

func NewEntitiesHandler(ctrl *base.HandlerController) *Handler {
	eb := entity_builder.NewEntityBuilder(ctrl.Logger, ctrl.DB, ctrl.CustomWriter)

	return &Handler{
		HandlerController: ctrl,
		entityBuilder:     eb,
	}
}
