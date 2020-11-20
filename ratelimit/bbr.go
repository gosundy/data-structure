package ratelimit

import (
	cpustat "data-struct/cpu"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var (
	cpu         int64
	decay       = 0.95
	cpuProcOnce sync.Once
)

type cpuGetter func() int64

// Op operations type.
type Op int

const (
	// Success opertion type: success
	Success Op = iota
	// Ignore opertion type: ignore
	Ignore
	// Drop opertion type: drop
	Drop
)

// DoneInfo done info.
type DoneInfo struct {
	Err error
	Op  Op
}

type BBR struct {
	passCounter     *RollingCounter
	rtCounter       *RollingCounter
	bucketDuration  time.Duration
	inFlight        int64
	winBucketPerSec int64
	cpu             cpuGetter
	prevDrop        atomic.Value
}

func NewBBR(window time.Duration, bucketCount int) *BBR {
	bbr := &BBR{}
	bbr.passCounter = NewRollingCounter(window, bucketCount)
	bbr.rtCounter = NewRollingCounter(window, bucketCount)
	bbr.bucketDuration = window / time.Duration(bucketCount)
	bbr.winBucketPerSec = int64(time.Second) / (int64(window) / int64(bucketCount))
	cpu := func() int64 {
		return atomic.LoadInt64(&cpu)
	}
	bbr.cpu = cpu
	cpuProcOnce.Do(startCPUProc)
	return bbr
}
func (bbr *BBR) Allow() (func(do DoneInfo), error) {
	if bbr.shouldDrop() {
		return nil, errors.New("limit exceed")
	}
	atomic.AddInt64(&bbr.inFlight, 1)
	stime := time.Now()
	return func(do DoneInfo) {
		rt := int64(time.Since(stime) / time.Millisecond)
		bbr.rtCounter.Add(float64(rt))
		atomic.AddInt64(&bbr.inFlight, -1)
		switch do.Op {
		case Success:
			bbr.passCounter.Add(1)
			return
		default:
			return
		}
	}, nil
}
func (bbr *BBR) shouldDrop() bool {
	fmt.Println(bbr.maxFlight())
	if bbr.cpu() < 800 {
		prevDrop, ok := bbr.prevDrop.Load().(time.Time)
		if !ok {
			return false
		}
		if time.Since(prevDrop) <= time.Second {
			inFlight := atomic.LoadInt64(&bbr.inFlight)
			return inFlight > 1 && inFlight > bbr.maxFlight()
		}
		return false
	}
	inFlight := atomic.LoadInt64(&bbr.inFlight)
	drop := inFlight > 1 && inFlight > bbr.maxFlight()
	if drop{
		bbr.prevDrop.Store(time.Now())
	}
	return drop
}
func (bbr *BBR) maxFlight() int64 {
	fmt.Println(bbr.winBucketPerSec,bbr.maxPassPerBucket(),bbr.minTtPerBucket())
	return int64(math.Floor(bbr.maxPassPerBucket()*float64(bbr.winBucketPerSec)*bbr.minTtPerBucket())/1000.0 + 0.5)
}
func (bbr *BBR) maxPassPerBucket() float64 {
	values := bbr.passCounter.Values()
	max := 0.0
	for _, value := range values {
		if value.SUM > max {
			max = value.SUM
		}
	}
	if max == 0 {
		max = 1
	}
	return max
}
func (bbr *BBR) minTtPerBucket() float64 {
	values := bbr.rtCounter.Values()
	min := math.MaxFloat64
	for _, value := range values {
		rt:=value.SUM/float64(value.COUNT)
		if  rt< min {
			min = rt
		}
	}
	return min
}
func startCPUProc() {
	cpustat.Init()
	go cpuproc()
}

// cpu = cpuᵗ⁻¹ * decay + cpuᵗ * (1 - decay)
func cpuproc() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("rate.limit.cpuproc() err(%+v)", err)
			go cpuproc()
		}
	}()
	ticker := time.NewTicker(time.Millisecond * 250)
	// EMA algorithm: https://blog.csdn.net/m0_38106113/article/details/81542863
	for range ticker.C {
		stat := &cpustat.Stat{}
		cpustat.ReadStat(stat)
		prevCPU := atomic.LoadInt64(&cpu)
		curCPU := int64(float64(prevCPU)*decay + float64(stat.Usage)*(1.0-decay))
		atomic.StoreInt64(&cpu, curCPU)
	}
}
