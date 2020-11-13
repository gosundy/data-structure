package lru

type Lru struct {
	//fast find data
	loc      map[int64]*Node
	listHead *Node
	listTail *Node
	//total capacity
	capacity int
	//current count
	curLength int
}
type Node struct {
	key  int64
	data interface{}
	pre  *Node
	next *Node
}

func NewLru(capacity int) *Lru {
	return &Lru{listHead: &Node{}, capacity: capacity, loc: make(map[int64]*Node)}
}
func (lru *Lru) Get(key int64) (interface{}, bool) {
	//if data exist, put it to first of list
	if node, ok := lru.loc[key]; ok {
		if lru.curLength == 1 {
			return node.data, false
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
		return node.data, false
	} else {
		return nil, true
	}

}
func (lru *Lru) Put(key int64, data interface{}) {
	if _, ok := lru.loc[key]; !ok {
		//add new node
		//if queue is full, let last node out
		node := &Node{data: data, key: key}
		node.next = lru.listHead.next
		lru.listHead.next = node
		node.pre = lru.listHead
		if node.next != nil {
			node.next.pre = node
		}
		lru.loc[key] = node
		lru.curLength += 1
		if lru.curLength == 1 {
			lru.listTail = node
		}
		if lru.curLength == lru.capacity+1 {
			deleteNode := lru.listTail
			lru.listTail = lru.listTail.pre
			delete(lru.loc, deleteNode.key)
			deleteNode = nil
			lru.curLength = lru.curLength - 1
		}
	}
}
