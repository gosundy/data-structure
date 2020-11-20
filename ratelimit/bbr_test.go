package ratelimit

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

//rt 为1s
//每秒 1次访问
func TestBBR(t *testing.T) {
	bbr := NewBBR(time.Second, 5)
	for i := 0; i < 100; i++ {
		doneFunc, err := bbr.Allow()
		if err != nil {
			t.Log(err)
		}
		time.Sleep(time.Second/10)
		doneFunc(DoneInfo{Op: Success})
	}
	t.Log(bbr.maxFlight())
}
func TestBBR2(t *testing.T) {

	limiter := NewBBR(time.Second * 5,50)
	var wg sync.WaitGroup
	var drop int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 300; i++ {
				f, err := limiter.Allow()
				if err != nil {
					atomic.AddInt64(&drop, 1)
				} else {
					count := rand.Intn(100)
					time.Sleep(time.Millisecond * time.Duration(count))
					f(DoneInfo{Op: Success})
				}
			}
		}()
	}
	wg.Wait()
	t.Logf("drop:%d ", drop)
}

