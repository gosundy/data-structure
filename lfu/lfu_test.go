package lfu

import (
	"log"
	"math/rand"
	"testing"
)

func TestLfu_Put01(t *testing.T) {
	lfu, _ := NewLfu(2)
	lfu.Put(1, "a")
	lfu.Put(2, "b")
	if len(lfu.heap.nodes) != 3 {
		t.Fatalf("expect:3, actual:%d", len(lfu.heap.nodes))
	}
	if lfu.heap.nodes[1].data.(string) != "a" {
		t.Fatalf("expect:a, actual:%v", lfu.heap.nodes[1].data)
	}
	if lfu.heap.nodes[2].data.(string) != "b" {
		t.Fatalf("expect:b, actual:%v", lfu.heap.nodes[2].data)
	}
}
func TestLfu_Put02(t *testing.T) {
	lfu, _ := NewLfu(2)
	lfu.Put(1, "a")
	lfu.Put(2, "b")
	lfu.Put(1, "a")
	if len(lfu.heap.nodes) != 3 {
		t.Fatalf("expect:3, actual:%d", len(lfu.heap.nodes))
	}
	if lfu.heap.nodes[1].data.(string) != "a" {
		t.Fatalf("expect:a, actual:%v", lfu.heap.nodes[1].data)
	}
	if lfu.heap.nodes[2].data.(string) != "b" {
		t.Fatalf("expect:b, actual:%v", lfu.heap.nodes[2].data)
	}
}

func TestLfu_Get01(t *testing.T) {
	lfu, _ := NewLfu(2)
	lfu.Put(1, "a")
	lfu.Put(2, "b")
	v, err := lfu.Get(1)
	if err != nil {
		log.Fatal(err)
	}
	if len(lfu.heap.nodes) != 3 {
		t.Fatalf("expect:3, actual:%d", len(lfu.heap.nodes))
	}
	if v.(string) != "a" {
		t.Fatalf("expect:a, actual:%v", v)
	}
	if lfu.heap.nodes[2].data.(string) != "a" {
		t.Fatalf("expect:a, actual:%v", lfu.heap.nodes[1].data)
	}
	if lfu.heap.nodes[1].data.(string) != "b" {
		t.Fatalf("expect:b, actual:%v", lfu.heap.nodes[2].data)
	}
}
func TestLfu_Get02(t *testing.T) {
	lfu, _ := NewLfu(2)
	lfu.Put(1, "a")
	lfu.Put(2, "b")
	v, err := lfu.Get(1)
	if err != nil {
		log.Fatal(err)
	}
	v, err = lfu.Get(2)
	if err != nil {
		log.Fatal(err)
	}
	v, err = lfu.Get(2)
	if err != nil {
		log.Fatal(err)
	}
	if len(lfu.heap.nodes) != 3 {
		t.Fatalf("expect:3, actual:%d", len(lfu.heap.nodes))
	}
	if v.(string) != "b" {
		t.Fatalf("expect:b, actual:%v", v)
	}
	if lfu.heap.nodes[1].data.(string) != "a" {
		t.Fatalf("expect:a, actual:%v", lfu.heap.nodes[1].data)
	}
	if lfu.heap.nodes[2].data.(string) != "b" {
		t.Fatalf("expect:b, actual:%v", lfu.heap.nodes[2].data)
	}
}
func BenchmarkLfu_Get(b *testing.B) {
	count := 300
	lfu, _ := NewLfu(300)
	for i := 0; i < count; i++ {
		lfu.Put(int64(i), i)
	}
	keys := rand.Perm(count)
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			data, err := lfu.Get(int64(key))
			if err != nil {
				log.Fatal(err)
			}
			if data.(int) != key {
				log.Fatalf("expect:%d, acutal:%d", key, data)
			}
		}
	}
	b.ReportAllocs()
}
