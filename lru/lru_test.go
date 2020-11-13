package lru

import (
	"testing"
)

//case:capacity=1
func TestLru_Get01(t *testing.T) {
	lru := NewLru(1)
	lru.Put(1, "1")
	_, isMiss := lru.Get(1)
	if isMiss == true {
		t.Fatalf("expect:hit, actual:miss")
	}
}
func TestLru_Get02(t *testing.T) {
	dataCount := 100
	lru := NewLru(dataCount)

	for i := 0; i < dataCount; i++ {
		if i%2 == 0 {
			lru.Put(int64(i), i)
		}
	}
	for i := 0; i < dataCount; i++ {
		_, isMiss := lru.Get(int64(i))
		if i%2 == 0 {
			if isMiss {
				t.Fatalf("expect:hit, actual:miss")
			}
		} else {
			if !isMiss {
				t.Fatalf("expect:miss, actual:hit")
			}
		}
	}
}
