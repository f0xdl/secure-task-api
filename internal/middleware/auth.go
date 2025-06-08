package middleware

import (
	"net/http"
	"sta/internal/handlers"
)

type BasicAuth struct {
	h *handlers.BasicAuthHandler
}

func NewBasicAuth(h *handlers.BasicAuthHandler) *BasicAuth {
	return &BasicAuth{h: h}
}

func (auth *BasicAuth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		if !ok {
			resp.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(resp, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if len(user) == 0 || len(pass) == 0 || !auth.h.Allow(user, pass) {
			resp.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(resp, "credentials are wrong", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(resp, req)
	})
}
