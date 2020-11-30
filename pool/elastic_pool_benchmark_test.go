// MIT License

// Copyright (c) 2018 Andy Pan

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package pool

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	RunTimes           = 1000000
	BenchParam         = 10
	BenchAntsSize      = 200000
	DefaultExpiredTime = 1000 * time.Millisecond
)

func demoFunc(ctx context.Context) {
	time.Sleep(BenchParam * time.Millisecond)
}

func demoPoolFunc(args interface{}) {
	n := args.(int)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func longRunningFunc() {
	for {
		runtime.Gosched()
	}
}

func longRunningPoolFunc(arg interface{}) {
	if ch, ok := arg.(chan struct{}); ok {
		<-ch
		return
	}
	for {
		runtime.Gosched()
	}
}

func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			go func() {
				demoFunc(context.Background())
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSemaphore(b *testing.B) {
	var wg sync.WaitGroup
	sema := make(chan struct{}, BenchAntsSize)

	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoFunc(context.Background())
				<-sema
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkAntsPool(b *testing.B) {
	var wg sync.WaitGroup
	p := NewPool(BenchAntsSize, DefaultExpiredTime)
	defer p.Close()

	b.StartTimer()
	//go func() {
	//	for {
	//		fmt.Println(p.workerRunCount,len(p.taskQueue))
	//		time.Sleep(time.Second)
	//
	//	}
	//}()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			p.Dispatch(func(ctx context.Context) {
				demoFunc(ctx)
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkGoroutinesThroughput(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			go demoFunc(context.Background())
		}
	}
}

func BenchmarkSemaphoreThroughput(b *testing.B) {
	sema := make(chan struct{}, BenchAntsSize)
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoFunc(context.Background())
				<-sema
			}()
		}
	}
}

func BenchmarkAntsPoolThroughput(b *testing.B) {
	p := NewPool(BenchAntsSize, DefaultExpiredTime)
	defer p.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			p.Dispatch(demoFunc)
		}
	}
	b.StopTimer()
}
