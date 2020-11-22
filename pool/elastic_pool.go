package pool

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
)

type Job struct {
	MaxQueueSize     int32
	CurrentQueueSize int32
	workerCh         chan *Worker
	ctx              context.Context
	cancel           context.CancelFunc
	workerRunCount   int32
}
type Worker struct {
	job *Job
}

func NewJob(currentQueueSize int32, maxQueueSize int32) *Job {
	job := &Job{}
	job.MaxQueueSize = maxQueueSize
	job.ctx, job.cancel = context.WithCancel(context.Background())
	job.CurrentQueueSize = currentQueueSize
	workerCh := make(chan *Worker, job.MaxQueueSize)
	for i := int32(0); i < job.MaxQueueSize; i++ {
		workerCh <- &Worker{job: job}
	}
	job.workerCh = workerCh
	return job
}
func (job *Job) Dispatch(task func(ctx context.Context) error) {
	worker := job.TakeWorker()
	atomic.AddInt32(&job.workerRunCount, 1)
	worker.Do(job.ctx, task)
}
func (job *Job) AddWorker(worker *Worker) {
	job.workerCh <- worker
}
func (job *Job) TakeWorker() *Worker {
	for {
		cur := atomic.LoadInt32(&job.CurrentQueueSize)
		curWorker := atomic.LoadInt32(&job.workerRunCount)
		if curWorker > cur {
			runtime.Gosched()
			continue
		}
		worker := <-job.workerCh
		return worker
	}

}
func (job *Job) Scale(ratio int32) {
	scaleSize := job.CurrentQueueSize * ratio
	if scaleSize > job.MaxQueueSize {
		scaleSize = job.MaxQueueSize
	}

	job.CurrentQueueSize = scaleSize
}
func (job *Job) DeScale(ratio int32) {
	scaleSize := job.CurrentQueueSize / ratio
	if scaleSize == 0 {
		scaleSize = 1
	}
	job.CurrentQueueSize = scaleSize

}
func (w *Worker) Do(ctx context.Context, task func(ctx context.Context) error) {
	go func() {
		defer atomic.AddInt32(&w.job.workerRunCount, -1)
		err := task(ctx)
		if err != nil {
			fmt.Printf("%v", err)
		}
		w.job.AddWorker(w)
	}()
}
