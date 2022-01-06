package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/k0mmsussert0d/fukaeri/internal/db"
	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"github.com/k0mmsussert0d/fukaeri/internal/workers"
)

func main() {
	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for range c {
			cancel()
		}
	}()

	db.InitCollections(ctx)
	workers.StartArchiving(ctx)
}
