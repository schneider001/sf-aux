package plugins

import (
	"log"

	"sf-aux/internal/models"
)

func Printer(in <-chan models.EventWithContext) {
	for {
		ev, ok := <-in
		if !ok {
			log.Println("Channel 'records' closed")
			return
		}

		bytes, err := ev.MarshalData()
		if err != nil {
			log.Println("ERROR", err)
			continue
		}

		log.Println(ev.Header, " <> ", string(bytes))
	}
}
