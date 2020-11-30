package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	QueueSize      int
	workerQueue    atomic.Value
	ctx            context.Context
	cancel         context.CancelFunc
	workerRunCount int64
	objectPool     *sync.Pool
	mu             sync.Mutex
	taskQueue      chan func(ctx context.Context)
	expireTime     time.Duration
	limitRate      int32
}
type Worker struct {
	pool *Pool
}

func NewPool(queueSize int, expireTime time.Duration) *Pool {
	pool := &Pool{}
	pool.QueueSize = queueSize
	pool.ctx, pool.cancel = context.WithCancel(context.Background())
	workerCh := make(chan *Worker, pool.QueueSize)
	_pool := sync.Pool{}
	_pool.New = func() interface{} {
		return &Worker{pool: pool}
	}
	pool.objectPool = &_pool
	for i := 0; i < pool.QueueSize; i++ {
		workerCh <- &Worker{pool: pool}
	}
	pool.workerQueue.Store(workerCh)
	pool.expireTime = expireTime
	pool.taskQueue = make(chan func(ctx context.Context), 0)
	return pool
}
func (pool *Pool) Dispatch(task func(ctx context.Context)) {
	for {
		//cold  start
		if pool.workerRunCount == 0 {
			worker := pool.TakeWorker()
			if worker != nil {
				worker.Do(pool.ctx)
				atomic.AddInt64(&pool.workerRunCount, 1)
			}
		}
		select {
		case pool.taskQueue <- task:
			return
		default:
		}
		//if task's queue is full, start one new worker
		worker := pool.TakeWorker()
		if worker != nil {
			worker.Do(pool.ctx)
			atomic.AddInt64(&pool.workerRunCount, 1)
		}
	}

}
func (job *Pool) AddWorker(worker *Worker) {
	workerQueue := job.workerQueue.Load().(chan *Worker)
	select {
	case workerQueue <- worker:
	default:
		job.objectPool.Put(worker)
	}
}
func (job *Pool) TakeWorker() *Worker {
	workerQueue := job.workerQueue.Load().(chan *Worker)
	select {
	case worker := <-workerQueue:
		return worker
	default:
		return nil
	}
}
func (job *Pool) Scale(scaleSize int) {
	job.mu.Lock()
	defer job.mu.Unlock()
	delta := scaleSize - job.QueueSize
	if delta <= 0 {
		return
	}
	workQueue := make(chan *Worker, scaleSize)

	for i := 0; i < delta; i++ {
		worker := job.objectPool.Get().(*Worker)
		workQueue <- worker
	}

	job.QueueSize = scaleSize
	job.workerQueue.Store(workQueue)
}

func (w *Worker) Do(ctx context.Context) {
	go func() {
		for task := range w.pool.taskQueue {
			if task == nil {
				return
			}
			task(ctx)
		}
	}()
}
func (pool *Pool) Close() {
	pool.cancel()
}
func (pool *Pool) Probe() {
	timer := time.NewTimer(time.Second * 2)
	go func() {
		for {
			select {
			case <-timer.C:
				startTime := time.Now()
				probe := make(chan struct{})
				pool.Dispatch(func(ctx context.Context) {
					probe <- struct{}{}
				})
				<-probe
				timer.Reset(time.Second * 2)
				if time.Since(startTime).Seconds() > 1 {
					pool.limitRate = 1
				} else {
					pool.limitRate = 0
				}
			case <-pool.ctx.Done():
				return
			}
		}

	}()

}
