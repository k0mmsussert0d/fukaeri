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
	// log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	conf.Init()
	log.Auto()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			<-c
			log.Warn().Print("RECEIVED EXIT SIGNAL! Terminating all workers")
			cancel()
		}
	}()

	db.InitCollections(ctx)
	workers.StartArchiving(ctx)
}
