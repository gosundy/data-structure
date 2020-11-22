package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Job struct {
	QueueSize      int
	workerQueue    atomic.Value
	ctx            context.Context
	cancel         context.CancelFunc
	workerRunCount int32
	objectPool     sync.Pool
	mu             sync.Mutex
}
type Worker struct {
	job *Job
}

func NewJob(queueSize int) *Job {
	job := &Job{}
	job.QueueSize = queueSize
	job.ctx, job.cancel = context.WithCancel(context.Background())
	workerCh := make(chan *Worker, job.QueueSize)
	pool := sync.Pool{}
	pool.New = func() interface{} {
		return &Worker{job: job}
	}
	job.objectPool = pool
	for i := 0; i < job.QueueSize; i++ {
		workerCh <- &Worker{job: job}
	}
	job.workerQueue.Store(workerCh)
	return job
}
func (job *Job) Dispatch(task func(ctx context.Context) error) {
	worker := job.TakeWorker()
	atomic.AddInt32(&job.workerRunCount, 1)
	worker.Do(job.ctx, task)
}
func (job *Job) AddWorker(worker *Worker) {
	workerQueue := job.workerQueue.Load().(chan *Worker)
	select {
	case workerQueue <- worker:
	default:
		job.objectPool.Put(worker)
	}
}
func (job *Job) TakeWorker() *Worker {
	workerQueue := job.workerQueue.Load().(chan *Worker)
	worker := <-workerQueue
	return worker
}
func (job *Job) Scale(queueSize int) {
	job.mu.Lock()
	defer job.mu.Unlock()
	scaleSize := queueSize
	if scaleSize < job.QueueSize {
		return
	}
	delta := scaleSize - job.QueueSize
	workQueue := make(chan *Worker, scaleSize)
	for i := 0; i < delta; i++ {
		worker := job.objectPool.Get().(*Worker)
		workQueue <- worker
	}
	job.QueueSize = scaleSize
	job.workerQueue.Store(workQueue)
}
func (job *Job) DeScale(queueSize int) {
	job.mu.Lock()
	defer job.mu.Unlock()
	scaleSize := queueSize
	if scaleSize > job.QueueSize {
		return
	}

	workQueue := make(chan *Worker, scaleSize)
	job.QueueSize = scaleSize
	job.workerQueue.Store(workQueue)

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
