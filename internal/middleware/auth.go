package middleware

import "net/http"

func Auth(h http.Handler) http.Handler {
	return h
}
