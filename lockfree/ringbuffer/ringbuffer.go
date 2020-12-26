package ringbuffer

import (
	"errors"
	"sync/atomic"
)

type Ring struct {
	front       int32
	rear        int32
	data        []interface{}
	len         int32
	rearStatus  int32
	frontStatus int32
}

var EmptyErr = errors.New("empty error")
var FullErr = errors.New("full error")

func NewRing(len int) *Ring {
	return &Ring{data: make([]interface{}, len), len: int32(len)}
}
func (ring *Ring) Put(data interface{}) error {
	for {
		if (ring.rear+1)%ring.len == ring.front {
			return FullErr
		}
		if atomic.CompareAndSwapInt32(&ring.rearStatus, 0, 0xff) {
			if (ring.rear+1)%ring.len == ring.front {
				ring.rearStatus = 0
				return FullErr
			}
			oldRear := ring.rear
			newRear := (oldRear + 1) % ring.len
			ring.data[oldRear] = data
			ring.rear = newRear
			ring.rearStatus = 0
			break
		}
	}
	return nil
}
func (ring *Ring) Get() (data interface{}, err error) {
	for {
		if ring.front == ring.rear {
			return nil, EmptyErr
		}
		if atomic.CompareAndSwapInt32(&ring.frontStatus, 0, 0xff) {
			if ring.front == ring.rear {
				ring.frontStatus = 0
				return nil, EmptyErr
			}
			oldFront := ring.front
			newFront := (oldFront + 1) % ring.len
			data = ring.data[oldFront]
			ring.front = newFront
			ring.frontStatus = 0
			break
		}
	}
	return data, nil

}
