package main

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/schneider001/sf-apis/go/sfgo"

	"sf-aux/internal/models"
	"sf-aux/internal/plugins"
)

func main() {
	logFile, err := os.OpenFile("/var/log/sf-aux/sf-aux.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	records := make(chan *sfgo.SysFlow, 16)
	events := make(chan models.EventWithContext, 16)

	socketPath, ok := os.LookupEnv("SF_SOCKET_PATH")
	if !ok {
		log.Fatal("SF_SOCKET_PATH must be set")
	}

	producer, err := plugins.MakeKafkaProducer("sf_events")
	if err != nil {
		log.Fatal("Make producer: ", err.Error())
	}

	agg := &plugins.AggregatorPlugin{}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer.Handle(events)
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
