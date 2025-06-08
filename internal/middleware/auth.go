package middleware

import (
	"errors"
	"net/http"
)

type IAuth interface {
	Allow(user, password string) bool
}

type BasicAuth struct {
	auth IAuth
}

func NewBasicAuth(auth IAuth) (*BasicAuth, error) {
	if auth == nil {
		return nil, errors.New("handler is nil")
	}
	return &BasicAuth{auth: auth}, nil
}

func (h *BasicAuth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		if !ok {
			resp.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(resp, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if len(user) == 0 || len(pass) == 0 || !h.auth.Allow(user, pass) {
			resp.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(resp, "credentials are wrong", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(resp, req)
	})
}
