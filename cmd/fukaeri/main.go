package main

import (
	"fmt"
	"strconv"

	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi"
)

func main() {
	chanapi.StartClient()
	gThreads := chanapi.Threads("g")

	for _, threadsOnPage := range gThreads {
		for _, threadInfo := range threadsOnPage.Threads {
			fmt.Print(chanapi.Thread("g", strconv.Itoa(threadInfo.No)))
		}
	}
}
