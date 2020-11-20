package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type RollingCounter struct {
	Window         time.Duration
	BucketCount    int
	Offset         int
	BucketDuration time.Duration
	buckets        []Bucket
	lastAccess     time.Time
	rwMutex        sync.RWMutex
}
type Bucket struct {
	SUM   float64
	COUNT int64
}

func NewRollingCounter(window time.Duration, bucketCount int) *RollingCounter {
	return &RollingCounter{Window: window, BucketCount: bucketCount, lastAccess: time.Now(), buckets: make([]Bucket, bucketCount), BucketDuration: window / time.Duration(bucketCount), rwMutex: sync.RWMutex{}}
}
func (metric *RollingCounter) Add(delta float64) {
	metric.rwMutex.Lock()
	defer metric.rwMutex.Unlock()
	spanDuration := time.Now().Sub(metric.lastAccess)
	span := timespan(spanDuration, metric.BucketDuration)
	offset := metric.Offset
	for i := 0; i < span && i < metric.BucketCount; i++ {
		offset = (offset + i + 1) % metric.BucketCount
		metric.buckets[offset].SUM = 0
		metric.buckets[offset].COUNT = 0
	}
	metric.buckets[offset].SUM += delta
	metric.buckets[offset].COUNT += 1
	metric.Offset = offset
	metric.lastAccess = metric.lastAccess.Add(time.Duration(span) * metric.BucketDuration)
}
func timespan(spanDuration time.Duration, durationBucketPer time.Duration) int {
	return int(spanDuration / durationBucketPer)
}
func (metric *RollingCounter) Values() []Bucket {
	return metric.buckets
}
func (metric *RollingCounter) String() string {
	return fmt.Sprintf("%v", metric.buckets)
}
