package middleware

import (
	"net/http"

	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	"github.com/oriiyx/fritz/app/core/utils/env"
)

func (am *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf := env.New()

		sessionID, err := r.Cookie(conf.Session.SessionCookieName)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		am.logger.Info().Str("session_id", sessionID.Value).Msg("session cookie")

		sessionDetails, err := am.session.ValidateSession(am.ctx, sessionID.Value)
		if err != nil {
			am.logger.Debug().Err(err).Str("session_id", sessionID.Value).Msg("Invalid session")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store token in context for downstream handlers
		ctx := ctxUtil.SetSession(r.Context(), sessionDetails)
		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
