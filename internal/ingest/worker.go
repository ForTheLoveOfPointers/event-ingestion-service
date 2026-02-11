package ingest

import (
	"context"
	"log"
	"sync"
)

type Worker struct {
	id int
}

func NewWorker(id int) *Worker {
	return &Worker{id: id}
}

func (w *Worker) work(b *Batcher, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker shutting down")
			return

		case <-b.FlushChannel:
			return
		}
	}
}
