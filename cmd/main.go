package main

import (
	"context"
	"github.com/f0xdl/secure-task-api/internal/handlers"
	"github.com/f0xdl/secure-task-api/internal/httpserver"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()

	log.Println("connecting Redis")
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("error connecting to Redis: %v", err)
	}

	log.Println("build authorization")
	auth := handlers.NewBasicAuthHandler()
	auth.AddCredentials(os.Getenv("AUTH_USERNAME"), os.Getenv("AUTH_PASSWORD"))

	log.Println("run webserver")
	err = httpserver.Run(ctx, rdb, auth, httpserver.Config{
		Host:      os.Getenv("HOST"),
		ApiPrefix: os.Getenv("API_PREFIX"),
	})
	if err != nil {
		log.Fatalf("error starting HTTP server: %v", err)
	}
}
