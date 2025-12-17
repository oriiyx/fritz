package ctx

import (
	"context"

	db "github.com/oriiyx/fritz/database/generated"
)

const sessionKey = key("session")

// GetSession Get session from context.
func GetSession(ctx context.Context) *db.Session {
	session, _ := ctx.Value(sessionKey).(*db.Session)
	return session
}

// SetSession Set session in context.
func SetSession(ctx context.Context, session *db.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}
