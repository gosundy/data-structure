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
	workerRunCount int32
	objectPool     *sync.Pool
	mu             sync.Mutex
	taskQueue      chan func(ctx context.Context)
	expireTime     time.Duration
	submitQueue    chan func(ctx context.Context)
	once           sync.Once
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
	pool.submitQueue = make(chan func(ctx context.Context), queueSize/3)
	pool.workerQueue.Store(workerCh)
	pool.expireTime = expireTime
	pool.taskQueue = make(chan func(ctx context.Context), queueSize*2)
	return pool
}
func (pool *Pool) Dispatch(task func(ctx context.Context)) {
	pool.once.Do(
		func() {
			for i := 0; i < 3; i++ {
				go func() {
					for {
					gettask:
						_task := <-pool.submitQueue
						for {
							if pool.workerRunCount == 0 {
								worker := pool.TakeWorker()
								if worker != nil {
									atomic.AddInt32(&pool.workerRunCount, 1)
									worker.Do(pool.ctx)
								}
							}
							select {
							case pool.taskQueue <- _task:
								goto gettask
							default:
							}
							if pool.workerRunCount == int32(pool.QueueSize) {
								continue
							}
							worker := pool.TakeWorker()
							if worker != nil {
								atomic.AddInt32(&pool.workerRunCount, 1)
								worker.Do(pool.ctx)
							}
						}
					}
				}()
			}

		})

	pool.submitQueue <- task

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
		timer.Stop()
		for {
			select {
			case task := <-w.pool.taskQueue:
				//timer.Stop()
				task(ctx)
				//timer.Reset(w.pool.expireTime)
			case <-timer.C:
				atomic.AddInt32(&w.pool.workerRunCount, -1)
				w.pool.AddWorker(w)
			case <-ctx.Done():
				atomic.AddInt32(&w.pool.workerRunCount, -1)
				return
			}
		}
	}()
}
func (pool *Pool) Close() {
	pool.cancel()
}
