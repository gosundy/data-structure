package ringbuffer

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRing01(t *testing.T) {
	t.Log("测试ring写，非并发")
	count := 10
	ring := NewRing(count + 1)
	for i := 0; i < count; i++ {
		err := ring.Put(i)
		if err != nil {
			t.Fatal("put with err:", err.Error())
		}
	}
	err := ring.Put(0)
	if err != FullErr {
		t.Fatal("should be full, actual not")
	}
}
func TestRing02(t *testing.T) {
	t.Log("测试ring写读，非并发")
	count := 10
	ring := NewRing(count + 1)
	for i := 0; i < count; i++ {
		err := ring.Put(i)
		if err != nil {
			t.Fatal("put with err:", err.Error())
		}
	}
	for i := 0; i < count; i++ {
		data, err := ring.Get()
		if err != nil {
			t.Fatal("get with err:", err.Error())
		}
		if data != i {
			t.Fatalf("should be %d, actual:%v", i, data)
		}
	}
}
func TestRing03(t *testing.T) {
	t.Log("测试ring写，一个读一个写")
	count := 100
	ring := NewRing(count + 1)
	readCon := 1
	writeCon := 1
	readCount := int32(0)
	wg := sync.WaitGroup{}
	wg.Add(readCon + writeCon)
	for i := 0; i < writeCon; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < count/writeCon; i++ {
				err := ring.Put(i)
				if err != nil {
					t.Fatal("put with err:", err.Error())
				}
			}
		}()
	}

	for i := 0; i < readCon; i++ {
		go func() {
			defer wg.Done()
			for {
				data, err := ring.Get()
				if err == nil {
					t.Log("read:", data)
					atomic.AddInt32(&readCount, 1)
				}
				if readCount == int32(count) {
					return
				}
			}

		}()

	}
	wg.Wait()
}

func TestRing05(t *testing.T) {
	t.Log("测试ring写，多个读多个写,ring<<写入的数据")

	ring := NewRing(10)
	readCon := 3
	writeCon := 3

	perWriteCount := int32(100000)
	curTotalWrite := int32(0)
	readCount := int32(0)
	wg := sync.WaitGroup{}
	wg.Add(readCon + writeCon)
	go func() {
		for {
			t.Log("cur read:", readCount, "cur write:",curTotalWrite)
			time.Sleep(time.Second)
		}
	}()
	defer func() {
		t.Log("cur read:", readCount, "cur write:",curTotalWrite)
	}()
	for i := 0; i < writeCon; i++ {
		go func(idx int) {
			defer wg.Done()
			k:=int32(0)
			for {

				err := ring.Put(0)
				if err == nil {
					atomic.AddInt32(&curTotalWrite, 1)
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
		t.Fatalf("expect empty, actual not")
	}
}
