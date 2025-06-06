package httpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

func Run(ctx context.Context, addr string, apiPrefix string) {
	apiHandler := NewHandler()
	rTask := apiHandler.RegisterRoutes()
	//rMetrics := apiHandler.RegisterMetrics()

	mux := http.NewServeMux()
	mux.Handle(apiPrefix+"/", http.StripPrefix(apiPrefix, rTask))

	log.Println("Listening on", addr)
	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println("Server error: ", err)
		}
	}()

	<-ctx.Done()
	log.Println("Stop serving on", addr)
	ctxStop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxStop); err != nil {
		log.Fatal("Server stopped error: ", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}

func healthHandler(resp http.ResponseWriter, req *http.Request) {
	log.Println("health handler")
	_, err := resp.Write([]byte("OK"))
	if err != nil {
		log.Fatal(err)
		return
	}
}
