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
func (job *Pool) Scale(queueSize int) {
	job.mu.Lock()
	defer job.mu.Unlock()
	scaleSize := queueSize
	delta := scaleSize - job.QueueSize
	if scaleSize <= 0 {
		scaleSize = 1
	}
	workQueue := make(chan *Worker, scaleSize)
	if delta > 0 {
		for i := 0; i < delta; i++ {
			worker := job.objectPool.Get().(*Worker)
			workQueue <- worker
		}
	}

	job.QueueSize = scaleSize
	job.workerQueue.Store(workQueue)
}

func (w *Worker) Do(ctx context.Context) {
	go func() {
		timer := time.NewTimer(w.pool.expireTime)
		for {
			select {
			case task := <-w.pool.taskQueue:
				task(ctx)
				timer.Reset(w.pool.expireTime)
			case <-timer.C:
				atomic.AddInt64(&w.pool.workerRunCount, -1)
				w.pool.AddWorker(w)
				return
			case <-ctx.Done():
				atomic.AddInt64(&w.pool.workerRunCount, -1)
				w.pool.AddWorker(w)
				return
			}
		}
	}()
}
func (pool *Pool) Close() {
	pool.cancel()
}
