package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/k0mmsussert0d/fukaeri/internal/conf"
	"github.com/k0mmsussert0d/fukaeri/internal/db"
	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"github.com/k0mmsussert0d/fukaeri/internal/workers"
)

func main() {
	conf.Init()

	log.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			<-c
			log.Logger().Warn("Received exit signal, terminating all workers")
			log.Logger().Sync()
			cancel()
		}
	}()

	db.InitCollections(ctx)
	workers.StartArchiving(ctx)
}
