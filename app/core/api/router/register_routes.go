package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/api/middleware"
	"github.com/oriiyx/fritz/app/core/api/middleware/requestlog"
	"github.com/oriiyx/fritz/app/core/services/auth"
	defHandler "github.com/oriiyx/fritz/app/core/services/definitions"
	"github.com/oriiyx/fritz/app/core/services/entities"
)

func (c *Controller) RegisterRoutes() {
	handlerFactory := base.NewHandlerControllerFactory(c.Logger, c.Queries, c.Validator, c.Kernel.Hooks(), c.DB, c.CustomWriter, c.Conf)
	authHandler := auth.NewAuth(c.Logger, c.Validator, c.Conf, c.DB, c.Store, c.Queries)
	am := middleware.NewAuthMiddleware(c.Logger, c.Queries, c.Conf, c.Ctx)

	c.Router.Route("/api/v1", func(r chi.Router) {
		r.Method(http.MethodPost, "/auth/login", requestlog.NewHandler(authHandler.Login, c.Logger))
		r.Method(http.MethodPost, "/auth/register", requestlog.NewHandler(authHandler.Register, c.Logger))

		// todo - finish implementation at a later time
		r.Method(http.MethodGet, "/auth/{provider}", requestlog.NewHandler(authHandler.ProviderLogin, c.Logger))
		r.Method(http.MethodGet, "/auth/{provider}/callback", requestlog.NewHandler(authHandler.ProviderCallback, c.Logger))
		r.Method(http.MethodGet, "/auth/{provider}/logout", requestlog.NewHandler(authHandler.ProviderLogout, c.Logger))

		definitionsHandler := defHandler.NewDefinitionsHandler(handlerFactory.Create("definitions"))
		r.Route("/definitions", func(definitions chi.Router) {
			definitions.Method(http.MethodGet, "/", requestlog.NewHandler(definitionsHandler.GetExisting, c.Logger))
			definitions.Method(http.MethodGet, "/data-component-types", requestlog.NewHandler(definitionsHandler.GetDataComponentTypes, c.Logger))
			definitions.Method(http.MethodPost, "/create", requestlog.NewHandler(definitionsHandler.Create, c.Logger))
		})

		entitiesHandler := entities.NewEntitiesHandler(handlerFactory.Create("entities"))
		r.Route("/entities", func(entities chi.Router) {
			entities.Method(http.MethodPost, "/{class_id}/read", requestlog.NewHandler(entitiesHandler.ReadEntity, c.Logger))
			entities.Method(http.MethodPost, "/{class_id}/create", requestlog.NewHandler(entitiesHandler.CreateEntity, c.Logger))
			entities.Method(http.MethodPost, "/{class_id}/update", requestlog.NewHandler(entitiesHandler.UpdateEntity, c.Logger))
			entities.Method(http.MethodPost, "/{class_id}/delete", requestlog.NewHandler(entitiesHandler.DeleteEntity, c.Logger))
		})

		r.Group(func(protectedRouter chi.Router) {
			protectedRouter.Use(am.AuthMiddleware)
			protectedRouter.Method(http.MethodGet, "/auth/me", requestlog.NewHandler(authHandler.MeHandler, c.Logger))
		})
	})
}
