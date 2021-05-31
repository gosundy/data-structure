package binary_heap

import "testing"

type TestNode struct {
	num int
}

func TestBinaryHeap(t *testing.T) {
	count:=100
	heap := NewBinaryHeap(count)
	for i := count; i > 0; i-- {
		err := heap.Push(&TestNode{num: i})
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 1; i <= count; i++ {
		popNode, err := heap.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if v, ok := popNode.(*TestNode); ok {
			if v.num != i {
				t.Fatalf("expected:%d,acutal:%d", i, v.num)
			}
		}
	}

}

func (n *TestNode) Less(b Compared) bool {
	if v, ok := b.(*TestNode); ok {
		if n.num < v.num {
			return true
		}
	}
	return false
}
