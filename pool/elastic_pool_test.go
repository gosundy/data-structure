package pool

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestJob(t *testing.T) {
	pool := NewPool(100000, time.Second*10)
	res := int32(0)
	count := 1000000
	start := time.Now()
	for i := 0; i < count; i++ {
		pool.Dispatch(func(ctx context.Context) {
			time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
			atomic.AddInt32(&res, 1)
			//if res%int32(10000) == 0 {
			//	pool.Scale(2*pool.QueueSize)
			//}
			//if res%int32(20000) == 0 {
			//	pool.Scale(pool.QueueSize / 2)
			//}
		})
	}
	for {
		if res == int32(count) {
			fmt.Println(res, time.Since(start), pool.QueueSize)
			return
		} else {
			mem := runtime.MemStats{}
			runtime.ReadMemStats(&mem)
			fmt.Println(mem.TotalAlloc, runtime.NumGoroutine())
			time.Sleep(time.Millisecond)
		}
	}

}
func TestJobForCpuCompute(t *testing.T) {
	count := 2
	res := uint64(0)
	workerCount := runtime.NumCPU()
	pool := NewPool(workerCount, time.Second)

	wg := sync.WaitGroup{}
	wg.Add(count)
	computeTask := func(ctx context.Context) {
		defer wg.Done()
		sum := uint64(0)
		for i := 0; i < 10000000; i++ {
			sum += uint64(i)
		}
		atomic.AddUint64(&res, sum)
	}
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(pool.workerRunCount)
		}
	}()
	for i := 0; i < count; i++ {
		pool.Dispatch(computeTask)
	}
	wg.Wait()
	t.Log(res)
}
func TestJobForIOCompute(t *testing.T) {
	count := 5000000
	workerCount := runtime.NumCPU()*40000
	pool := NewPool(workerCount, time.Second)

	wg := sync.WaitGroup{}
	wg.Add(count)
	computeTask := func(ctx context.Context) {
		defer wg.Done()
		time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
		//time.Sleep(10 * time.Millisecond)
	}
	go func() {
		for {
			//if pool.workerRunCount == int64(pool.QueueSize) {
			//	if pool.workerRunCount > 10000 {
			//		pool.Scale(int(pool.workerRunCount + 1000))
			//	} else {
			//		pool.Scale(int(pool.workerRunCount * 2))
			//	}
			//
			//}
			fmt.Println(atomic.LoadInt64(&pool.workerRunCount), len(pool.taskQueue))
			time.Sleep(time.Second)
		}
	}()
	pool.Probe()
	go func() {
		for i := 0; i < count; i++ {
			pool.Dispatch(computeTask)
		}
	}()
	defer pool.Close()

	wg.Wait()
}
func TestJobForIOComputeWithDeathLock(t *testing.T) {
	count := runtime.NumCPU() + 1
	workerCount := runtime.NumCPU()
	pool := NewPool(workerCount, time.Second*10)

	wg := sync.WaitGroup{}
	wg.Add(count)
	res := make(chan int)

	computeTask := func(ctx context.Context) {
		defer wg.Done()
		res <- 1
	}

	for i := 0; i < count; i++ {
		pool.Dispatch(computeTask)
	}
	go func() {
		for i := 0; i < count; i++ {
			<-res
		}

	}()
	wg.Wait()
}
