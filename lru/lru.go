package lru

import (
	"unsafe"
)

//defined for fast find data in map
type LruData interface {
	HashCode() int64
}
type Lru struct {
	//fast find data
	loc      map[int64]unsafe.Pointer
	listHead *Node
	listTail *Node
	//total capacity
	capacity int
	//current count
	curLength int
}
type Node struct {
	data LruData
	pre  *Node
	next *Node
}

func NewLru(capacity int) *Lru {
	return &Lru{listHead: &Node{}, capacity: capacity, loc: make(map[int64]unsafe.Pointer)}
}
func (lru *Lru) Get(data LruData) (LruData, bool, error) {
	hashCode := data.HashCode()
	//if data exist, put it to first of list
	if nodeAddressPointer, ok := lru.loc[hashCode]; ok {
		node := (*Node)(nodeAddressPointer)
		if lru.curLength == 1 {
			return node.data, false, nil
		}
		if node.next != nil {
			node.next.pre = node.pre
			node.pre.next = node.next
		} else {
			lru.listTail = node.pre
		}
		node.next = lru.listHead.next
		node.pre = lru.listHead
		lru.listHead.next = node
		node.next.pre = node
		return node.data, false, nil
	} else {
		//add new node
		//if queue is full, let last node out
		node := &Node{data: data}
		node.next = lru.listHead.next
		lru.listHead.next = node
		node.pre = lru.listHead
		if node.next != nil {
			node.next.pre = node
		}
		lru.loc[data.HashCode()] = unsafe.Pointer(node)
		lru.curLength += 1
		if lru.curLength == 1 {
			lru.listTail = node
		}
		if lru.curLength == lru.capacity+1 {
			deleteNode := lru.listTail
			lru.listTail = lru.listTail.pre
			delete(lru.loc, deleteNode.data.HashCode())
			deleteNode = nil
			lru.curLength = lru.curLength - 1
		}
		return node.data, true, nil
	}
}
