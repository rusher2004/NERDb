package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rusher2004/nerdb/engine"
	"github.com/rusher2004/nerdb/listener"
)

func main() {
	ctx := context.Background()
	cl := http.Client{Timeout: 20 * time.Second}

	if err := engine.RunKillmails(ctx, &cl); err != nil {
		log.Fatalf("error running killmails: %v", err)
	}

	listener.Listen(ctx, &cl)
}
