package middleware

import (
	"net/http"

	"github.com/rs/xid"

	ctxUtil "github.com/oriiyx/fritz/utils/ctx"
)

const requestIDHeaderKey = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(requestIDHeaderKey)
		if requestID == "" {
			requestID = xid.New().String()
		}

		ctx = ctxUtil.SetRequestID(ctx, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
