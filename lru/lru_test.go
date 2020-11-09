package lru

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

type Data struct {
	idx int64
}

func (data Data) HashCode() int64 {
	return data.idx
}

//case:capacity=1
func TestLru_Get01(t *testing.T) {
	lru := NewLru(1)
	data := Data{
		idx: 1,
	}
	_, _, err := lru.Get(data)
	if err != nil {
		t.Fatal(err)
	}
	_data, isMiss, err := lru.Get(data)
	if err != nil {
		t.Fatal(err)
	}
	if isMiss == true {
		t.Fatalf("expect:hit, actual:miss")
	}
	data0 := _data.(Data)
	if data.idx != data0.idx {
		t.Fatalf("expect:equal, actual:not equal")
	}
}

//case:capacity=2
func TestLru_Get02(t *testing.T) {
	lru := NewLru(2)
	data01 := Data{idx: 1}
	data02 := Data{idx: 2}
	data03 := Data{idx: 3}

	_, isMiss, err := lru.Get(data01)
	if err != nil {
		t.Fatal(err)
	}
	if isMiss == false {
		t.Fatalf("expect:miss, actual:hit")
	}
	_, isMiss, err = lru.Get(data02)
	if err != nil {
		t.Fatal(err)
	}
	if isMiss == false {
		t.Fatalf("expect:miss, actual:hit")
	}
	_, isMiss, err = lru.Get(data03)
	if err != nil {
		t.Fatal(err)
	}
	if isMiss == false {
		t.Fatalf("expect:miss, actual:hit")
	}

	_data01I, isMiss, err := lru.Get(data01)
	if err != nil {
		t.Fatal(err)
	}
	if isMiss == false {
		t.Fatalf("expect:miss, actual:hit")
	}
	_data01 := _data01I.(Data)
	if _data01.idx != data01.idx {
		t.Fatalf("expect:equal, actual:not equal")
	}
}

//test add data cost time
func TestLru_Get03(t *testing.T) {
	datas := make([]Data, 100)

	ids := rand.Perm(100)
	for _, id := range ids {
		datas = append(datas, Data{idx: int64(id)})
	}
	lru := NewLru(100)
	start := time.Now()
	for _, data := range datas {
		_, _, err := lru.Get(data)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("put 100 count cost time:%s", time.Since(start))
}
func BenchmarkLru_Get(b *testing.B) {
	datas := make([]Data, 300)
	for i := 0; i < 3; i++ {
		ids := rand.Perm(100)
		for _, id := range ids {
			datas = append(datas, Data{idx: int64(id)})
		}
	}
	lru := NewLru(100)
	missCount := int32(0)
	for i := 0; i < b.N; i++ {
		for _, data := range datas {
			_, isMiss, err := lru.Get(data)
			if err != nil {
				b.Fatal(err)
			}
			if isMiss {
				atomic.AddInt32(&missCount, 1)
			}
		}

	}
	b.ReportAllocs()
	b.Logf("miss count:%d, current Len:%d", missCount, lru.curLength)
}
