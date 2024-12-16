package plugins

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/schneider001/sf-apis/go/sfgo"

	"sf-aux/internal/models"
)

type AggregatorPlugin struct {
	header sfgo.SFHeader
}

func stringifyType(t sfgo.RecUnionTypeEnum) (string, error) {
	switch t {
	case sfgo.SF_CONT:
		return "SF_CONT", nil
	case sfgo.SF_PROCESS:
		return "SF_PROCESS", nil
	case sfgo.SF_FILE:
		return "SF_FILE", nil
	case sfgo.SF_PROC_EVT:
		return "SF_PROC_EVT", nil
	case sfgo.SF_NET_FLOW:
		return "SF_NET_FLOW", nil
	case sfgo.SF_FILE_FLOW:
		return "SF_FILE_FLOW", nil
	case sfgo.SF_FILE_EVT:
		return "SF_FILE_EVT", nil
	case sfgo.SF_PROC_FLOW:
		return "SF_PROC_FLOW", nil
	case sfgo.SF_NET_EVT:
		return "SF_NET_EVT", nil
	default:
		return "", errors.New("unknown type")
	}
}

func (a *AggregatorPlugin) Handle(in <-chan *sfgo.SysFlow, out chan<- models.EventWithContext) {
	timeFile, err := os.OpenFile("/sysflow/sf-aux/check_perf/time_aggregator", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Open time file: ", err)
	}
	defer timeFile.Close()

	for {
		start := time.Now()

		rec, ok := <-in
		if !ok {
			log.Println("Channel 'in' closed")
			return
		}

		switch rec.Rec.UnionType {
		case sfgo.SF_HEADER:
			a.header = *rec.Rec.SFHeader
		default:
			typ, err := stringifyType(rec.Rec.UnionType)
			if err != nil {
				log.Printf("Error determining record type: %v", err)
				continue
			}

			out <- models.EventWithContext{
				Header: a.header,
				Type:   typ,
				Data:   *rec.Rec,
			}
		}

		duration := time.Since(start).Microseconds()
		_, _ = fmt.Fprintf(timeFile, "%d\n", duration)
	}
}
