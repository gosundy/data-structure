package ratelimit

import (
	"testing"
	"time"
)

func TestFunnel_Add(t *testing.T) {
	rateLimit := NewFunnel(1000, 1000)
	for i := 0; i < 1000; i++ {
		rateLimit.Add()
	}
	start := time.Now()
	for {
		if rateLimit.current != 0 {
			rateLimit.Take()
		} else {
			break
		}
	}
	t.Logf("cost:%d us, expect about: %d us", time.Since(start).Microseconds(), time.Second/1000)
}
