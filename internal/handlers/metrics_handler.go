package handlers

import (
	"log"
	"net/http"
)

func MetricsHandler(resp http.ResponseWriter, _ *http.Request) {
	_, err := resp.Write([]byte("1")) //mock
	if err != nil {
		log.Fatal(err)
		return
	}
}
