package ingest

import (
	"context"
	"sync"
)

type WorkerPool struct {
	numWorkers int
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{numWorkers: numWorkers}
}

// TODO: Add a WaitGroup
func (w *WorkerPool) Start(b *Batcher, ctx context.Context) {

	for i := 0; i < w.numWorkers; i++ {
		worker := NewWorker(i)
		w.wg.Add(1)
		go worker.work(b, ctx, &w.wg)
	}

}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
}
