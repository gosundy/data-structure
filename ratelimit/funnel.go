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
	datas    []interface{}
}

func NewFunnel(maxCapacity int64, perDataDuration time.Duration) *Funnel {
	return &Funnel{rwMutex: sync.RWMutex{}, capacity: maxCapacity, current: 0, leakRate: perDataDuration, lastTake: time.Now().Add(-perDataDuration), datas: make([]interface{}, maxCapacity)}
}
func (f *Funnel) Add(data interface{}) bool {
	if f.current == f.capacity {
		return false
	}
	if f.current > f.capacity {
		panic(fmt.Sprintf("funnel's current data count:%d more than max capacity:%d", f.current, f.capacity))
	}
	f.rwMutex.Lock()
	defer f.rwMutex.Unlock()
	f.datas[f.current] = data
	f.current++
	return true
}
func (f *Funnel) Take() (interface{}, bool) {
	if f.current == 0 {
		return nil, false
	}
	if f.current < 0 {
		panic(fmt.Sprintf("funnel's current data count:%d less than 0", f.current))
	}
	if time.Now().Sub(f.lastTake).Microseconds()-f.leakRate.Microseconds() < 0 {
		return nil, false
	}
	f.rwMutex.Lock()
	defer f.rwMutex.Unlock()
	f.current--
	f.lastTake = time.Now()
	return f.datas[f.current], true
}
