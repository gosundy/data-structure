package linklist

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkLink(b *testing.B) {
	link := NewLink()
	writeCount := 10000
	readCon := 10
	writeCon := 10
	b.ReportAllocs()
	for i:=0;i<b.N;i++ {
		wg := sync.WaitGroup{}
		readCount := int32(0)
		wg.Add(readCon + writeCon)
		//写协程
		for i := 0; i < writeCon; i++ {
			go func() {
				defer wg.Done()
				for k := 0; k < writeCount/writeCon; k++ {
					_ = link.Put(&Node{data: k})
				}
			}()
		}
		//读协程
		for i := 0; i < readCon; i++ {
			go func() {
				defer wg.Done()
				for {
					_, err := link.Get()
					if err == nil {
						atomic.AddInt32(&readCount, 1)
					}
					if readCount == int32(writeCount) {
						return
					}
				}
			}()
		}
		wg.Wait()
		_, err := link.Get()
		if err != EmptyErr {
			b.Fatal("expect empty, actual not empty")
		}
	}
}
