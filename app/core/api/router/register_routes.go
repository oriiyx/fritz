package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/api/middleware/requestlog"
	defHandler "github.com/oriiyx/fritz/app/core/services/definitions"
	"github.com/oriiyx/fritz/app/core/services/entities"
)

func (c *Controller) RegisterRoutes() {
	handlerFactory := base.NewHandlerControllerFactory(c.Logger, c.Queries, c.Validator, c.Kernel.Hooks())

	c.Router.Route("/api/v1", func(r chi.Router) {
		definitionsHandler := defHandler.NewDefinitionsHandler(handlerFactory.Create("definitions"))
		r.Route("/definitions", func(definitions chi.Router) {
			definitions.Method(http.MethodGet, "/data-component-types", requestlog.NewHandler(definitionsHandler.GetDataComponentTypes, c.Logger))
		})

		entitiesHandler := entities.NewEntitiesHandler(handlerFactory.Create("entities"))
		r.Route("/entities", func(entities chi.Router) {
			entities.Method(http.MethodPost, "/create", requestlog.NewHandler(entitiesHandler.CreateEntity, c.Logger))
		})
	})
}
