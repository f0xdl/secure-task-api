package httpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sta/internal/middleware"
	"sta/internal/middleware/ratelimit"
	"time"
)

func Run(ctx context.Context, addr string, apiPrefix string) {
	apiHandler := NewHandler()
	rateLimit := ratelimit.NewRateIpLimit(2, 0.2)
	rTask := middleware.Logger(rateLimit.Handle(apiHandler.RegisterRoutes()))
	rMetrics := middleware.Logger(middleware.Auth(rateLimit.Handle(apiHandler.RegisterMetrics())))
	mux := http.NewServeMux()
	mux.Handle(apiPrefix+"/", http.StripPrefix(apiPrefix, rTask))
	adminPrefix := apiPrefix + "/admin"
	mux.Handle(adminPrefix+"/", http.StripPrefix(adminPrefix, rMetrics))

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
