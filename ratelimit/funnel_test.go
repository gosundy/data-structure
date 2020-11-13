package ratelimit

import (
	"testing"
	"time"
)

func TestFunnel_Add(t *testing.T) {
	rateLimit:=NewFunnel(100,10*time.Millisecond)
	for i:=0;i<100;i++{
		rateLimit.Add(i)
	}
	start:=time.Now()
	for {
		if rateLimit.current!=0{
			rateLimit.Take()
		}else{
			break
		}
	}
	t.Logf("cost:%d us, expect about: %d us",time.Now().Sub(start).Microseconds(),100*10*time.Millisecond/1000)
}