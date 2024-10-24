package models

import (
	"github.com/schneider001/sf-apis/go/sfgo"
)

type EventWithContext struct {
	Header sfgo.SFHeader `json:"header"`

	Type string        `json:"type"`
	Data sfgo.RecUnion `json:"data"`
}

func (e *EventWithContext) MarshalData() ([]byte, error) {
	return e.Data.MarshalJSON()
}
