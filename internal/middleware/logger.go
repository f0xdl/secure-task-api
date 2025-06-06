package middleware

import "net/http"

func Logger(h http.Handler) http.Handler {
	return h
}
