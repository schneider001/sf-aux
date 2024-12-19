package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
	"time"

	"sf-aux/internal/models"
	"sf-aux/internal/plugins"

	"github.com/schneider001/sf-apis/go/sfgo"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.SetBlockProfileRate(1)

	logFile, err := os.OpenFile("/var/log/sf-aux/sf-aux.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	go func() {
		log.Println("Starting pprof server on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	records := make(chan *sfgo.SysFlow, 16000)
	events := make(chan models.EventWithContext, 16000)

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
		monitorChannels(records, events, "/sysflow/sf-aux/check_perf/records_fill", "/sysflow/sf-aux/check_perf/events_fill")
	}()

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

func monitorChannels(records chan *sfgo.SysFlow, events chan models.EventWithContext, recordsLogPath, eventsLogPath string) {
	recordsLogFile, err := os.OpenFile(recordsLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open records log file: %v", err)
	}
	defer recordsLogFile.Close()

	eventsLogFile, err := os.OpenFile(eventsLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open events log file: %v", err)
	}
	defer eventsLogFile.Close()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		recordsLen := len(records)
		eventsLen := len(events)

		recordsLine := fmt.Sprintf("%d\n", recordsLen)
		eventsLine := fmt.Sprintf("%d\n", eventsLen)

		if _, err := recordsLogFile.WriteString(recordsLine); err != nil {
			log.Printf("Failed to write to records log: %v", err)
		}
		if _, err := eventsLogFile.WriteString(eventsLine); err != nil {
			log.Printf("Failed to write to events log: %v", err)
		}
	}
}
