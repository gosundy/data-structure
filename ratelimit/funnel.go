package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type Funnel struct {
	rwMutex  sync.RWMutex
	capacity int64
	current  int64
	leakRate time.Duration
	lastTake time.Time
}

func NewFunnel(maxCapacity int64, qps uint64) *Funnel {
	leakRate := time.Second / time.Duration(qps)
	return &Funnel{rwMutex: sync.RWMutex{}, capacity: maxCapacity, current: 0, leakRate: leakRate, lastTake: time.Now().Add(-leakRate)}
}
func (f *Funnel) Add() bool {
	if f.current == f.capacity {
		return false
	}
	if f.current > f.capacity {
		panic(fmt.Sprintf("funnel's current data count:%d more than max capacity:%d", f.current, f.capacity))
	}
	f.rwMutex.Lock()
	defer f.rwMutex.Unlock()
	f.current++
	return true
}
func (f *Funnel) Take() bool {
	if f.current == 0 {
		return false
	}
	if f.current < 0 {
		panic(fmt.Sprintf("funnel's current data count:%d less than 0", f.current))
	}
	if time.Now().Sub(f.lastTake).Microseconds()-f.leakRate.Microseconds() < 0 {
		return false
	}
	f.rwMutex.Lock()
	defer f.rwMutex.Unlock()
	f.current--
	f.lastTake = time.Now()
	return true
}
