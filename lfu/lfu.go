package lfu

import (
	"errors"
	"fmt"
)

type Lfu struct {
	nodeIndex map[int64]*BinaryHeapNode
	capacity  int
	heap      *BinaryHeap
}
type BinaryHeap struct {
	curLen   int
	capacity int
	nodes    []*BinaryHeapNode
}
type BinaryHeapNode struct {
	key       int64
	frequency int
	data      interface{}
	heapIdx   int
}

func NewLfu(capacity int) (*Lfu, error) {
	if capacity < 1 {
		return nil, errors.New("capacity must more than 0")
	}
	return &Lfu{capacity: capacity, nodeIndex: make(map[int64]*BinaryHeapNode), heap: &BinaryHeap{curLen: 1, capacity: capacity, nodes: make([]*BinaryHeapNode, capacity+1)}}, nil
}
func (lfu *Lfu) Put(key int64, data interface{}) {
	if _, ok := lfu.nodeIndex[key]; ok {
		return
	}
	addedNode := lfu.heap.Put(key, data)
	lfu.nodeIndex[key] = addedNode
}
func (lfu *Lfu) Get(key int64) (interface{}, error) {
	if node, ok := lfu.nodeIndex[key]; ok {
		node.frequency += 1
		lfu.heap.shift(node.heapIdx)
		return node.data, nil
	} else {
		return nil, errors.New("key not found")
	}
}
func (heap *BinaryHeap) Put(key int64, data interface{}) *BinaryHeapNode {
	//pop top
	if heap.curLen > heap.capacity && heap.capacity >= 1 {
		heap.nodes[1].key = key
		heap.nodes[1].data = data
		return heap.nodes[1]
	} else {
		addedNode := &BinaryHeapNode{key: key, data: data, frequency: 1, heapIdx: heap.curLen}
		heap.nodes[heap.curLen] = addedNode
		//shift node  from low to high
		parentIdx := heap.curLen / 2
		needComparedIdx := heap.curLen
		for parentIdx > 0 && parentIdx != 1 {
			if heap.nodes[needComparedIdx].frequency < heap.nodes[parentIdx].frequency {
				//swap
				tmp := heap.nodes[needComparedIdx]
				heap.nodes[needComparedIdx] = heap.nodes[parentIdx]
				heap.nodes[needComparedIdx].heapIdx = needComparedIdx
				heap.nodes[parentIdx] = tmp
				heap.nodes[parentIdx].heapIdx = parentIdx
			}
			parentIdx = parentIdx / 2
			needComparedIdx = parentIdx
		}
		heap.curLen += 1
		return addedNode
	}
}
func (heap *BinaryHeap) shift(nodeIdx int) {
	//shift node from high to low
	minIdx := nodeIdx
	for {
		oldMinIdx := minIdx
		leftNodeIdx := minIdx * 2
		rightNodeIdx := minIdx*2 + 1
		if leftNodeIdx <= heap.curLen {
			if heap.nodes[leftNodeIdx].frequency < heap.nodes[minIdx].frequency {
				minIdx = leftNodeIdx
			}
		}
		if rightNodeIdx <= heap.curLen-1 {
			if heap.nodes[rightNodeIdx].frequency < heap.nodes[minIdx].frequency {
				minIdx = rightNodeIdx
			}
		}
		//swap
		if minIdx != oldMinIdx {
			tmp := heap.nodes[minIdx]
			heap.nodes[minIdx] = heap.nodes[oldMinIdx]
			heap.nodes[minIdx].heapIdx = minIdx
			heap.nodes[oldMinIdx] = tmp
			heap.nodes[oldMinIdx].heapIdx = oldMinIdx
		} else {
			break
		}
	}
}
func (lfu Lfu) String() string {
	strs := make([]string, len(lfu.heap.nodes)-1)
	for idx, node := range lfu.heap.nodes[1:] {
		strs[idx] = fmt.Sprintf("%v:%v", node.key, node.data)
	}
	return fmt.Sprintf("%v", strs)
}
