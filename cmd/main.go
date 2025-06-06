package main

import (
	"context"
	"os/signal"
	httpserver "sta/internal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()
	//logger
	//redis
	httpserver.Run(ctx, ":8080", "/api/v1")
}
