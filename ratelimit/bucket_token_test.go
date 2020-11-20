package ratelimit

import (
	"testing"
	"time"
)

func TestNewBucketToken(t *testing.T) {
	bucketToken := NewBucketToken(10, 1000)
	bucketToken.Start()
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second / 1000)
		token := bucketToken.Take()
		t.Log(token)
	}
	bucketToken.Close()

}
