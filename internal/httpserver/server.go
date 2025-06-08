package httpserver

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"sta/internal/handlers"
	"sta/internal/middleware"
	"sta/internal/middleware/ratelimit"
	"time"
)

func Run(ctx context.Context, rdb *redis.Client, addr string, apiPrefix string) error {
	log.Println("run webserver")
	hRecovery := middleware.Recover
	apiHandler := NewHandler()
	authHandler := handlers.NewBasicAuthHandler()
	authHandler.AddCredentials("admin", "T3st") //ONLY FOR MOCK
	auth, err := middleware.NewBasicAuth(authHandler)
	if err != nil {
		return err
	}

	var limit ratelimit.LimitHandler
	if rdb == nil {
		limit, err = ratelimit.NewRateIpLimit(2, 0.2)
	} else {
		limit, err = ratelimit.NewRedisRts(rdb, 1, 2, 5*time.Second)
	}
	if err != nil {
		return err
	}

	log.Println("configure routes")
	taskRouter := middleware.Logger(hRecovery(limit.Handle(apiHandler.RegisterRoutes())))
	metricRouter := middleware.Logger(hRecovery(limit.Handle(auth.Handle(apiHandler.RegisterMetrics()))))
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
	return nil
}
