package requestlog

import (
	"net/http"
)

type NoLogHandler struct {
	handler http.Handler
}

func NewNoLogHandler(h http.HandlerFunc) *NoLogHandler {
	return &NoLogHandler{
		handler: h,
	}
}

func (h *NoLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r2 := new(http.Request)
	*r2 = *r
	rcc := &readCounterCloser{r: r.Body}
	r2.Body = rcc
	w2 := &responseStats{w: w}

	h.handler.ServeHTTP(w2, r2)
}
