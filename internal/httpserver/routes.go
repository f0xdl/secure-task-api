package httpserver

import (
	"net/http"
	"sta/internal/handlers"
)

type Routes struct{}

func NewRoutes() *Routes {
	return &Routes{}
}

func (h *Routes) ApiRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /status", handlers.HealthHandler)
	r.HandleFunc("POST /task", handlers.TaskHandler)
	return r
}

func (h *Routes) AdminMetrics() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /metrics", handlers.MetricsHandler)
	return r
}
