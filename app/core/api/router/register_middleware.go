package router

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	internalMiddleware "github.com/oriiyx/fritz/app/core/api/middleware"
	"github.com/oriiyx/fritz/app/core/kernel/hooks"
)

func (c *Controller) RegisterUses() {
	err := c.Kernel.Hooks().Trigger(c.Ctx, hooks.HookBeforeRouterMiddleware, c.Router)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("Could not trigger hook: %s", hooks.HookBeforeRouterMiddleware)
	}

	c.Router.Use(middleware.Logger)
	c.Router.Use(internalMiddleware.RequestID)
	c.Router.Use(internalMiddleware.JSONMiddleware)

	// CORS Origins
	var allowedOrigins []string
	allowedOrigins = []string{
		"http://localhost:*",
	}
	if len(c.Conf.Server.CORSOrigins) != 0 {
		allowedOrigins = c.Conf.Server.CORSOrigins
	}

	// CORS Methods
	var allowedMethods []string
	allowedMethods = []string{
		"GET", "POST", "PUT", "DELETE", "OPTIONS",
	}
	if len(c.Conf.Server.CORSMethods) != 0 {
		allowedMethods = c.Conf.Server.CORSMethods
	}

	// CORS Headers
	var allowedHeaders []string
	allowedHeaders = []string{
		"Accept", "Authorization", "Content-Type", "X-CSRF-Token",
	}
	if len(c.Conf.Server.CORSHeaders) != 0 {
		allowedHeaders = c.Conf.Server.CORSHeaders
	}

	// CORS Exposed Headers
	var exposedHeaders []string
	exposedHeaders = []string{
		"Link",
	}
	if len(c.Conf.Server.CORSExposedHeaders) != 0 {
		exposedHeaders = c.Conf.Server.CORSExposedHeaders
	}

	// CORS Allow Credentials
	allowedCredentials := true
	if c.Conf.Server.CORSAllowCredentials != nil {
		allowedCredentials = *c.Conf.Server.CORSAllowCredentials
	}

	// CORS Max Age
	maxAge := 300
	if c.Conf.Server.CORSMaxAge != nil {
		maxAge = *c.Conf.Server.CORSMaxAge
	}

	c.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   exposedHeaders,
		AllowCredentials: allowedCredentials, // Set to true if you need cookies for cross-site requests
		MaxAge:           maxAge,             // Maximum value not to check for CORS again for a certain duration
	}))

	err = c.Kernel.Hooks().Trigger(c.Ctx, hooks.HookAfterRouterMiddleware, c.Router)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("Could not trigger hook: %s", hooks.HookAfterRouterMiddleware)
	}

	c.Logger.Info().Str("environment", c.Conf.Server.ENV).Msg("Environment")
}
