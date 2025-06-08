package middleware

import (
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//log.Printf("panic recovered: %v\n%s", err, debug.Stack())
				http.Error(resp, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(resp, req)
	})
}
