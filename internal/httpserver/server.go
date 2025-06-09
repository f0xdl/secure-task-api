package httpserver

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"sta/internal/middleware"
	"sta/internal/middleware/ratelimit"
	"time"
)

type Config struct {
	Host      string
	ApiPrefix string
}

func Run(ctx context.Context, rdb *redis.Client, auth middleware.IAuth, cfg Config) error {
	hRecovery := middleware.Recover
	apiRoutes := NewRoutes()
	aa, err := middleware.NewBasicAuth(auth)
	if err != nil {
		return err
	}

	var limit ratelimit.LimitHandler
	if rdb == nil {
		log.Println("Using in-memory rate limiter")
		limit, err = ratelimit.NewRateIpLimit(2, 0.2)
	} else {
		log.Println("Using Redis rate limiter")
		limit, err = ratelimit.NewRedisRts(rdb, 1, 2, 5*time.Second)
	}
	if err != nil {
		return err
	}

	log.Println("configure routes")
	taskRoutes := middleware.Logger(hRecovery(limit.Handle(apiRoutes.ApiRoutes())))
	metricRoutes := middleware.Logger(hRecovery(limit.Handle(aa.Handle(apiRoutes.AdminMetrics()))))
	mux := http.NewServeMux()
	mux.Handle(cfg.ApiPrefix+"/", http.StripPrefix(cfg.ApiPrefix, taskRoutes))
	mux.Handle(cfg.ApiPrefix+"/admin/", http.StripPrefix(cfg.ApiPrefix+"/admin", metricRoutes))

	log.Println("Listening on", cfg.Host)
	server := &http.Server{Addr: cfg.Host, Handler: mux}
	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println("Server error: ", err)
		}
	}()
	<-ctx.Done()

	log.Println("Stop serving on", cfg.Host)
	ctxStop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxStop); err != nil {
		log.Fatal("Server stopped error: ", err)
	} else {
		log.Println("Server gracefully stopped")
	}
	return nil
}
