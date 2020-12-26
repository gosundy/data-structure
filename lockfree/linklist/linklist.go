package linklist

import (
	"errors"
	"sync"
)

var EmptyErr = errors.New("empty")

type Link struct {
	head    *Node
	tail    *Node
	headMux sync.Mutex
	tailMux sync.Mutex
}
type Node struct {
	data interface{}
	next *Node
}

func NewLink() *Link {
	link := &Link{}
	node := new(Node)
	link.head = node
	link.tail = node
	return link
}
func (link *Link) Put(node *Node) error {
	link.tailMux.Lock()
	link.tail.next = node
	link.tail = link.tail.next
	link.tailMux.Unlock()
	return nil
}
func (link *Link) Get() (node *Node, err error) {
	if link.head.next == nil {
		return nil, EmptyErr
	}
	link.headMux.Lock()
	if link.head.next == nil {
		link.headMux.Unlock()
		return nil, EmptyErr
	}
	node = link.head.next
	link.head = link.head.next
	link.headMux.Unlock()
	return node, nil
}
