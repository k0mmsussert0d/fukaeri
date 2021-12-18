package main

import (
	"context"

	"github.com/k0mmsussert0d/fukaeri/internal/db"
	"github.com/k0mmsussert0d/fukaeri/internal/workers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db.InitCollections(ctx)
	workers.StartArchiving(ctx)
}
