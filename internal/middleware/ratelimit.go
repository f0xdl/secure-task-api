package middleware

import "net/http"

func RateLimit(h http.Handler) http.Handler {
	return h
}
