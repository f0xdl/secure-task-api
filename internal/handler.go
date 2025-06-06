package httpserver

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /status", healthHandler)
	r.HandleFunc("POST /task", healthHandler) //TODO: task handler
	return r
}

func (h *Handler) RegisterMetrics() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /metrics", healthHandler) //TODO: task handler
	return r
}
