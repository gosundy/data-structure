package pool

import (
	"context"
	"testing"
	"time"
)

func TestJob(t *testing.T) {
	job := NewJob(10, 200)
	for i := 0; i < 10000; i++ {
		job.Dispatch(func(ctx context.Context) error {
			//t.Log("task complete")
			return nil
		})
	}
	t.Log("task complete")
	time.Sleep(time.Second * 10)
}
