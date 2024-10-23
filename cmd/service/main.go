package main

import (
	"sync"

	"sf-aux/internal/plugins"

	"github.com/schneider001/sf-apis/go/sfgo"
)

const (
	socketPath = "/sock/sysflow.sock"
)

func main() {
	records := make(chan *sfgo.SysFlow, 16)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		plugins.Printer(records)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		plugins.Reader(socketPath, records)
	}()

	wg.Wait()
}
