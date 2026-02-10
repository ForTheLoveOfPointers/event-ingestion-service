package ingest

import (
	"context"
	"event-ingestion-service/internal/httpserver/types"
	"time"
)

type Batcher struct {
	EventBatch    []types.Event
	maxSize       int
	FlushInterval time.Duration
	FlushChannel  chan []types.Event
	flushTicker   *time.Ticker
}

/*
The batcher runs in memory, that is to say we are using heap memory (RAM) for all of this.
As such, I would not recommend doing big batch sizes, and would go for maybe a batch size
of 500-1000 and no more than that when running dockerized.
*/
func Make(maxSize int, flushInterval time.Duration) *Batcher {
	return &Batcher{
		EventBatch:    make([]types.Event, maxSize),
		maxSize:       maxSize,
		FlushInterval: flushInterval,
		FlushChannel:  make(chan []types.Event, maxSize),
	}
}

func (b *Batcher) Start(ingestor *Ingest, ctx context.Context) {

	b.flushTicker = time.NewTicker(b.FlushInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				b.flush()
				close(b.FlushChannel)
				return
			case <-b.flushTicker.C:
				if len(b.EventBatch) > 0 {
					b.flush()
				}
			case e := <-ingestor.ch:
				b.add(e)
			}
		}
	}()
}

func (b *Batcher) add(e types.Event) {
	b.EventBatch = append(b.EventBatch, e)
	if len(b.EventBatch) >= b.maxSize {
		b.flush()
	}
}

func (b *Batcher) flush() {
	defer b.flushTicker.Reset(b.FlushInterval)

	batchCopy := make([]types.Event, len(b.EventBatch))
	copy(batchCopy, b.EventBatch)
	b.FlushChannel <- batchCopy
	b.EventBatch = b.EventBatch[:0]
}
