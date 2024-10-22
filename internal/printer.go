package internal

import (
	"log"

	"github.com/schneider001/sf-apis/go/sfgo"
)

func Handle(records <-chan *sfgo.SysFlow) {
	for {
		rec, ok := <-records
		if !ok {
			log.Println("Channel 'records' closed")
			return
		}

		switch rec.Rec.UnionType {
		case sfgo.SF_HEADER:
			log.Printf("[SF_HEADER] %v\n", rec.Rec.SFHeader)
		case sfgo.SF_CONT:
			log.Printf("[SF_CONT] %v\n", rec.Rec.Container)
		case sfgo.SF_PROCESS:
			log.Printf("[SF_PROCESS] %v\n", rec.Rec.Process)
		case sfgo.SF_FILE:
			log.Printf("[SF_FILE] %v\n", rec.Rec.FileEvent)
		case sfgo.SF_PROC_EVT:
			log.Printf("[SF_PROC_EVT] %v\n", rec.Rec.ProcessEvent)
		case sfgo.SF_NET_FLOW:
			log.Printf("[SF_NET_FLOW] %v\n", rec.Rec.NetworkFlow)
		case sfgo.SF_FILE_FLOW:
			log.Printf("[SF_FILE_FLOW] %v\n", rec.Rec.FileFlow)
		case sfgo.SF_FILE_EVT:
			log.Printf("[SF_FILE_EVT] %v\n", rec.Rec.FileEvent)
		case sfgo.SF_PROC_FLOW:
			log.Printf("[SF_PROC_FLOW] %v\n", rec.Rec.ProcessFlow)
		case sfgo.SF_NET_EVT:
			log.Printf("[SF_NET_EVT] %v\n", rec.Rec.NetworkEvent)
		default:
			log.Println("Error unsupported SysFlow Type: ", rec.Rec.UnionType)
		}
	}
}