package middleware

import (
	"log"
	"net/http"
	"time"
)

type ResponseWriterWithStatus struct {
	http.ResponseWriter
	Status int
}

func NewResponseWriterWithStatus(w http.ResponseWriter) *ResponseWriterWithStatus {
	return &ResponseWriterWithStatus{w, http.StatusOK}
}

func (w *ResponseWriterWithStatus) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()
		respStat := NewResponseWriterWithStatus(resp)
		next.ServeHTTP(respStat, req)
		dur := time.Since(start)
		log.Printf("%s %d %s HTTP (%s)", req.Method, respStat.Status, req.URL.String(), dur)
	})
}
