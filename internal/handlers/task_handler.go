package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// TaskHandler {"value":0}
func TaskHandler(resp http.ResponseWriter, req *http.Request) {
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
