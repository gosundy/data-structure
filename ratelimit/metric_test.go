package ratelimit

import (
	"testing"
	"time"
)

func TestNewRollingCounter(t *testing.T) {
	count := 10
	counter := NewRollingCounter(time.Second, count)
	for i := 0; i < count; i++ {
		counter.Add(1)
		time.Sleep(time.Second / time.Duration(count))
	}
	for i := 0; i < count; i++ {
		if counter.buckets[i].SUM != 1 {
			t.Fatalf("expect:1, actual:%f", counter.buckets[i].SUM)
		}
	}
}
