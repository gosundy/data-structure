package ringbuffer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkRing(b *testing.B) {
	ring := NewRing(10)
	readCon := 1
	writeCon := 1
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		perWriteCount := int32(1000)
		readCount := int32(0)
		wg := sync.WaitGroup{}
		wg.Add(readCon + writeCon)
		for i := 0; i < writeCon; i++ {
			go func(idx int) {
				defer wg.Done()
				k := int32(0)
				for {
					err := ring.Put(0)
					if err == nil {
						k++
					}
					if k == perWriteCount {
						//t.Log("writer:",idx,"exit")
						return
					}

				}
			}(i)
		}

		for i := 0; i < readCon; i++ {
			go func() {
				defer wg.Done()
				for {
					_, err := ring.Get()
					if err == nil {
						atomic.AddInt32(&readCount, 1)
					}
					if readCount == perWriteCount*int32(writeCon) {
						return
					}

				}

			}()

		}
		wg.Wait()
		_, err := ring.Get()
		if err != EmptyErr {
			b.Fatalf("expect empty, actual not")
		}
	}
}
