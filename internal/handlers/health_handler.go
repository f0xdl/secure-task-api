package handlers

import (
	"log"
	"net/http"
)

func HealthHandler(resp http.ResponseWriter, req *http.Request) {
	_, err := resp.Write([]byte("OK"))
	if err != nil {
		log.Fatal(err)
		return
	}
}
