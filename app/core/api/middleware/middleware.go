package middleware

import (
	"context"

	"github.com/oriiyx/fritz/app/core/services/session"
	"github.com/oriiyx/fritz/app/core/utils/env"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

type AuthMiddleware struct {
	ctx     context.Context
	logger  *zerolog.Logger
	session *session.Manager
	queries *db.Queries
	envConf *env.Conf
}

func NewAuthMiddleware(
	logger *zerolog.Logger,
	queries *db.Queries,
	envConf *env.Conf,
	ctx context.Context,
) *AuthMiddleware {
	sessionManager := session.NewManager(envConf, queries, logger)
	return &AuthMiddleware{
		ctx:     ctx,
		logger:  logger,
		queries: queries,
		envConf: envConf,
		session: sessionManager,
	}
}
