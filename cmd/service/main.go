package main

import (
	"log"
	"os"
	"sync"

	"github.com/schneider001/sf-apis/go/sfgo"

	"sf-aux/internal/models"
	"sf-aux/internal/plugins"
)

func main() {
	records := make(chan *sfgo.SysFlow, 16)
	events := make(chan models.EventWithContext, 16)

	socketPath, ok := os.LookupEnv("SF_SOCKET_PATH")
	if !ok {
		log.Fatal("SF_SOCKET_PATH must be set")
	}

	agg := &plugins.AggregatorPlugin{}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		plugins.Printer(events)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		agg.Handle(records, events)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		plugins.Reader(socketPath, records)
	}()

	wg.Wait()
}
