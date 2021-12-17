package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := apiclient.New(ctx)
	gThreads := client.Threads("g")

	for _, threadsOnPage := range gThreads {
		for _, threadInfo := range threadsOnPage.Threads {
			fmt.Print(client.Thread("g", strconv.Itoa(threadInfo.No)))
		}
	}
}
