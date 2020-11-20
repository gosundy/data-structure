package pool

import (
	"context"
	"fmt"
	"time"
)

type Job struct {
	WorkerQueueSize   int
	MaxQueueCount     int
	CurrentQueueCount int
	workers           []chan *Worker
	ctx               context.Context
	cancel            context.CancelFunc
}
type Worker struct {
	job *Job
}

func NewJob(workerQueueSize int, maxQueueCount int) *Job {
	job := &Job{}
	job.WorkerQueueSize = workerQueueSize
	job.MaxQueueCount = maxQueueCount
	job.CurrentQueueCount = maxQueueCount
	job.workers = make([]chan *Worker, maxQueueCount)
	for i := 0; i < maxQueueCount; i++ {
		job.workers[i] = make(chan *Worker, workerQueueSize)
		for k:=0;k<job.WorkerQueueSize;k++{
			job.workers[i]<-&Worker{job: job}
		}
	}
	job.ctx, job.cancel = context.WithCancel(context.Background())

	return job
}
func (job *Job) Dispatch(task func(ctx context.Context) error) {
	worker := job.TakeWorker()
	worker.Do(job.ctx, task)
}
func (job *Job) AddWorker(worker *Worker) {
	for _, workerCh := range job.workers {
		if len(workerCh) < job.WorkerQueueSize {
			workerCh <- worker
			break
		}
	}
}
func (job *Job) TakeWorker() *Worker {
	for {
		select {
		case <-job.ctx.Done():
			return nil
		default:

		}
		for i := 0; i < job.CurrentQueueCount; i++ {
			if len(job.workers[i]) > 0 {
				worker := <-job.workers[i]
				return worker
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

}

func (w *Worker) Do(ctx context.Context, task func(ctx context.Context) error) {
	go func() {
		err := task(ctx)
		if err != nil {
			fmt.Printf("%v", err)
		}
		w.job.AddWorker(w)
	}()
}
