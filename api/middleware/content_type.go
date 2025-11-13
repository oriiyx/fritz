package middleware

import (
	"net/http"
)

const (
	HeaderKeyContentType       = "Content-Type"
	HeaderValueContentTypeJSON = "application/json;charset=utf8"
)

// JSONMiddleware sets the Content-Type to application/json for all responses.
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderKeyContentType, HeaderValueContentTypeJSON)
		next.ServeHTTP(w, r)
	})
}
