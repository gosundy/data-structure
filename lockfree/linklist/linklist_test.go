package linklist

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestLink(t *testing.T) {
	link := NewLink()
	for i := 0; i < 10; i++ {
		_ = link.Put(&Node{data: i})
	}
	for i := 0; i < 10; i++ {
		node, _ := link.Get()
		if node.data != i {
			t.Fatal("expect equal, actual not equal")
		}
	}
}

//测试多协程，一个读，一个写
func TestLinkMutiCon01(t *testing.T) {
	t.Log("测试多协程，一个读，一个写")
	link := NewLink()
	count := 100000
	wg := sync.WaitGroup{}
	wg.Add(2)
	//写协程
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			_ = link.Put(&Node{data: i})
		}
	}()
	//读协程
	go func() {
		defer wg.Done()
		i := 0
		for {
			node, err := link.Get()
			if err == nil {
				if node.data != i {
					t.Fatal("expect equal, actual not equal")
				}
				i++
			}
			if i == count {
				return
			}
		}
	}()
	wg.Wait()
}

//测试多协程，多个读，一个写
func TestLinkMutiCon02(t *testing.T) {
	t.Log("测试多协程，多个读，一个写")
	link := NewLink()
	writeCount := 1000000
	wg := sync.WaitGroup{}
	readCon := 10
	readCount := int32(0)
	wg.Add(readCon + 1)
	//写协程
	go func() {
		defer wg.Done()
		for i := 0; i < writeCount; i++ {
			_ = link.Put(&Node{data: i})
		}
	}()
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
}

//测试多协程，多个读，多个写
func TestLinkMutiCon03(t *testing.T) {
	t.Log("测试多协程，多个读，多个写")
	link := NewLink()
	writeCount := 1000000
	wg := sync.WaitGroup{}
	readCon := 10
	writeCon := 10
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
		t.Fatal("expect empty, actual not empty")
	}
}
