package workerpool

import (
	"context"
	"cryptopricealerter/internal/repository"
	"sync"
)

type WorkerPool struct {
	JobChan     chan Job
	Workers     []*Worker
	WorkerCount int
	AlertRepo   repository.AlertRepository
	Ctx         context.Context
	WG          *sync.WaitGroup
	Cancel      context.CancelFunc
}

func NewWorkerPool(chanSize int, workerCount int, repo repository.AlertRepository, ctx context.Context, cancel context.CancelFunc) *WorkerPool {
	return &WorkerPool{
		JobChan:     make(chan Job, chanSize),
		Workers:     make([]*Worker, workerCount),
		WorkerCount: workerCount,
		AlertRepo:   repo,
		Ctx:         ctx,
		WG:          &sync.WaitGroup{},
		Cancel: cancel,
	}
}

func (wp *WorkerPool) Start() {
	for i := range wp.Workers {
		wp.Workers[i] = NewWorker(uint(i), wp.JobChan, wp.AlertRepo, wp.Ctx, wp.WG)
		wp.WG.Add(1)
		go wp.Workers[i].Start()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.JobChan)
	wp.Cancel()
	wp.WG.Wait()
}
