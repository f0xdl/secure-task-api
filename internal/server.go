package httpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sta/internal/handlers"
	"sta/internal/middleware"
	"sta/internal/middleware/ratelimit"
	"time"
)

func Run(ctx context.Context, addr string, apiPrefix string) {
	log.Println("run webserver")
	apiHandler := NewHandler()
	authHandler := handlers.NewBasicAuthHandler()
	authHandler.AddCredentials("admin", "T3st") //ONLY FOR MOCK
	auth := middleware.NewBasicAuth(authHandler)
	limit := ratelimit.NewRateIpLimit(2, 0.2)
	r := middleware.Recover

	log.Println("configure routes")
	taskRouter := middleware.Logger(r(limit.Handle(apiHandler.RegisterRoutes())))
	metricRouter := middleware.Logger(r(auth.Handle(limit.Handle(apiHandler.RegisterMetrics()))))
	mux := http.NewServeMux()
	mux.Handle(apiPrefix+"/", http.StripPrefix(apiPrefix, taskRouter))
	adminPrefix := apiPrefix + "/admin"
	mux.Handle(adminPrefix+"/", http.StripPrefix(adminPrefix, metricRouter))

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
