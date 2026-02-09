package ingest

import (
	"event-ingestion-service/internal/httpserver/types"
)

type Ingest struct {
	ch chan types.Event
}

func New(bufferSize int) *Ingest {
	return &Ingest{ch: make(chan types.Event, bufferSize)}
}

func (i *Ingest) Enqueue(e types.Event) bool {
	select {
	case i.ch <- e:
		return true
	default:
		return false
	}
}
