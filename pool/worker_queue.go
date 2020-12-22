package pool

import (
	"fmt"
)

type WorkerQueue struct {
	size   int
	worker []*Worker
	rear   int
	front  int
}

func NewWorkerQueue(size int) *WorkerQueue {
	if size <= 0 {
		panic(fmt.Errorf("size %d must more than zero", size))
	}
	queue := &WorkerQueue{
		size:   size,
		worker: make([]*Worker, size),
	}
	return queue
}
func Dequeue()(){

}
