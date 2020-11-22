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
	job := NewJob(10, 100000)
	res := int32(0)
	count := 1000000
	start := time.Now()
	for i := 0; i < count; i++ {
		job.Dispatch(func(ctx context.Context) error {
			time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
			atomic.AddInt32(&res, 1)
			if res%int32(100) == 0 {
				job.Scale(2)
			}
			return nil
		})
	}
	for {
		if res == int32(count) {
			fmt.Println(res, time.Since(start), job.CurrentQueueSize)
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
	count := 100
	res := uint64(0)
	workerCount := int32(runtime.NumCPU())
	job := NewJob(workerCount, workerCount)

	wg := sync.WaitGroup{}
	wg.Add(count)
	computeTask := func(ctx context.Context) error {
		defer wg.Done()
		sum := uint64(0)
		for i := 0; i < 1000000000; i++ {
			sum += uint64(i)
		}
		atomic.AddUint64(&res, sum)
		return nil
	}

	for i := 0; i < count; i++ {
		job.Dispatch(computeTask)
	}
	wg.Wait()
	t.Log(res)
}
func TestJobForIOCompute(t *testing.T) {
	count := 10000
	workerCount := int32(runtime.NumCPU()) * 100
	job := NewJob(workerCount, workerCount)

	wg := sync.WaitGroup{}
	wg.Add(count)
	computeTask := func(ctx context.Context) error {
		defer wg.Done()

		time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
		return nil
	}

	for i := 0; i < count; i++ {
		job.Dispatch(computeTask)
	}
	wg.Wait()
}
func TestJobForIOComputeWithDeathLock(t *testing.T) {
	count := runtime.NumCPU() + 1
	workerCount := int32(runtime.NumCPU())
	job := NewJob(workerCount, workerCount)

	wg := sync.WaitGroup{}
	wg.Add(count)
	res := make(chan int)

	computeTask := func(ctx context.Context) error {
		defer wg.Done()
		res <- 1
		return nil
	}

	for i := 0; i < count; i++ {
		job.Dispatch(computeTask)
	}
	go func() {
		for i := 0; i < count; i++ {
			<-res
		}

	}()
	wg.Wait()
}
