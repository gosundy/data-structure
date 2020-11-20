package ratelimit

import (
	"sync"
	"time"
)

type BucketToken struct {
	capacity uint64
	//一个token耗费的时间
	rate time.Duration
	//当前的token数量
	current uint64
	mutex   sync.Mutex
	done    chan struct{}
	once    sync.Once
}

func NewBucketToken(maxCapacity uint64, qps uint64) *BucketToken {
	return &BucketToken{capacity: maxCapacity, rate: time.Second / time.Duration(qps), current: 0, mutex: sync.Mutex{}, done: make(chan struct{}), once: sync.Once{}}
}
func (rateLimit *BucketToken) Take() bool {
	if rateLimit.current == 0 {
		return false
	}
	if rateLimit.current < 0 {
		panic("not reach")
	}
	rateLimit.mutex.Lock()
	defer rateLimit.mutex.Unlock()
	rateLimit.current--
	return true
}
func (rateLimit *BucketToken) Close() {
	rateLimit.done <- struct{}{}
	rateLimit.once.Do(func() {
		close(rateLimit.done)
	})

}
func (rateLimit *BucketToken) Start() {
	go rateLimit.generate()
}
func (rateLimit *BucketToken) generate() {
	ticker := time.NewTicker(rateLimit.rate)
	for {
		select {
		case <-ticker.C:
			rateLimit.mutex.Lock()
			if rateLimit.current == rateLimit.capacity {
				rateLimit.mutex.Unlock()
				continue
			}
			if rateLimit.current > rateLimit.capacity {
				rateLimit.mutex.Unlock()
				panic("no reach")
			}
			rateLimit.current++
			rateLimit.mutex.Unlock()
		case <-rateLimit.done:
			ticker.Stop()
			return

		}
	}

}
