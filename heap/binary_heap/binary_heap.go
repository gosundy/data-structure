package binary_heap

import "errors"

var ErrHeapIsFull = errors.New("heap is full")
var ErrHeapIsEmpty = errors.New("heap is empty")

type Compared interface {
	Less(data Compared) bool
}
type BinaryHeap struct {
	Capacity int
	Length   int
	nodes    []Compared
}

func NewBinaryHeap(size int) *BinaryHeap {
	return &BinaryHeap{Capacity: size, Length: 0, nodes: make([]Compared, size+1)}
}

func (b *BinaryHeap) Push(insertNode Compared) error {
	if b.Length > b.Capacity {
		return ErrHeapIsFull
	}
	//append data
	b.Length = b.Length + 1
	b.nodes[b.Length] = insertNode
	//shift data
	insertNodeIdx := b.Length
	parentIdx := insertNodeIdx / 2
	for parentIdx != 0 {
		if !b.nodes[insertNodeIdx].Less(b.nodes[parentIdx]) {
			break
		}
		tmpNode := b.nodes[parentIdx]
		b.nodes[parentIdx] = b.nodes[insertNodeIdx]
		b.nodes[insertNodeIdx] = tmpNode
		insertNodeIdx = parentIdx
		parentIdx = insertNodeIdx / 2
	}
	return nil
}

func (b *BinaryHeap) Pop() (Compared, error) {
	if b.Length == 0 {
		return nil, ErrHeapIsEmpty
	}
	if b.Length == 1 {
		b.Length = 0
		return b.nodes[1], nil
	}
	popNode := b.nodes[1]
	b.nodes[1] = b.nodes[b.Length]
	b.Length = b.Length - 1
	//shift
	needShiftNodeIdx := 1
	for {
		// compared with left Node and right node
		leftNodeIdx := needShiftNodeIdx * 2
		rightNodeIdx := needShiftNodeIdx*2 + 1

		if leftNodeIdx <= b.Length && rightNodeIdx <= b.Length {
			if !b.nodes[leftNodeIdx].Less(b.nodes[needShiftNodeIdx]) && !b.nodes[rightNodeIdx].Less(b.nodes[needShiftNodeIdx]) {
				break
			}
			var minNodeIdx int
			if b.nodes[leftNodeIdx].Less(b.nodes[rightNodeIdx]) {
				minNodeIdx = leftNodeIdx
			} else {
				minNodeIdx = rightNodeIdx
			}
			tempNode := b.nodes[minNodeIdx]
			b.nodes[minNodeIdx] = b.nodes[needShiftNodeIdx]
			b.nodes[needShiftNodeIdx] = tempNode
			needShiftNodeIdx = minNodeIdx
			continue
		}
		if leftNodeIdx <= b.Length {
			if b.nodes[leftNodeIdx].Less(b.nodes[needShiftNodeIdx]) {
				tempNode := b.nodes[leftNodeIdx]
				b.nodes[leftNodeIdx] = b.nodes[needShiftNodeIdx]
				b.nodes[needShiftNodeIdx] = tempNode
				needShiftNodeIdx = leftNodeIdx
			}
		}
		break
	}

	return popNode, nil
}
