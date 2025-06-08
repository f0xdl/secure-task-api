package httpserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /status", healthHandler)
	r.HandleFunc("POST /task", taskHandler)
	return r
}

func (h *Handler) RegisterMetrics() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /metrics", healthHandler) //TODO: metrics handler
	return r
}

// taskHandler {"value":0}
func taskHandler(resp http.ResponseWriter, req *http.Request) {
	//parsing
	body := req.Body
	defer body.Close()
	raw, err := io.ReadAll(body)
	if err != nil {
		http.Error(resp, "wrong data", http.StatusUnprocessableEntity)
		return
	}
	data := map[string]int{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		http.Error(resp, "wrong JSON format", http.StatusUnprocessableEntity)
		return
	}
	v, ok := data["value"]
	if !ok {
		http.Error(resp, "no value", http.StatusUnprocessableEntity)
		return
	}
	//calculate
	b, _ := json.Marshal(map[string]int{"result": v * v})
	_, err = resp.Write(b)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
}

func healthHandler(resp http.ResponseWriter, req *http.Request) {
	_, err := resp.Write([]byte("OK"))
	if err != nil {
		log.Fatal(err)
		return
	}
}
