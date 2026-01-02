package workerpool

import (
	"context"
	"cryptopricealerter/internal/repository"
	"fmt"
	"log"
	"sync"
)

type Worker struct {
	id uint
	jobChan <-chan Job
	alertRepo repository.AlertRepository
	ctx context.Context
	wg *sync.WaitGroup
}

func NewWorker(id uint, jobchan <-chan Job, repo repository.AlertRepository, ctx context.Context, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id: id,
		jobChan: jobchan,
		alertRepo: repo,
		ctx: ctx,
		wg: wg,
	}
}

func (w *Worker) Start() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		case job, ok := <-w.jobChan:
			if !ok {
				return
			}
			switch job.Condition {
			case ">":
				if job.ActualPrice > job.Price {
					handleTriggered(w.alertRepo, job, w.id)
				}
			case "<":
				if job.ActualPrice < job.Price {
					handleTriggered(w.alertRepo, job, w.id)
				}
			}

		}

	}
}

func handleTriggered(alertRepo repository.AlertRepository, job Job, workerID uint) {
	if err := alertRepo.MarkTriggered(job.ID); err != nil {
		log.Println("Error:", err)
		return
	}
	fmt.Println("Алерт", job.ID, "сработал:", job.ActualPrice, job.Condition, job.Price, "Воркер:", workerID)
}