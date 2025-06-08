package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os/signal"
	"sta/internal/httpserver"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalln(err)
	}
	err = rdb.FlushDB(ctx).Err()
	if err != nil {
		log.Fatalln(err)
	}
	err = httpserver.Run(ctx, rdb, ":8080", "/api/v1")
	if err != nil {
		log.Fatalln(err)
	}
}
