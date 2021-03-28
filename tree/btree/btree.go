package btree

import "math"

type BTree struct {
	root *BTreeNode
	rank int
}
type BTreeNode struct {
	keys   []int
	keyLen int
	parent *BTreeNode
	childs []*BTreeNode
}

func NewBtree(rank int) *BTree {
	return &BTree{rank: rank}
}
func (tree *BTree) Insert(data int) {
	if tree.root == nil {
		tree.root = &BTreeNode{keys: make([]int, tree.rank), keyLen: 0, childs: make([]*BTreeNode, tree.rank+1)}
		tree.root.keys[0] = data
		tree.root.keyLen = 1
		return
	}
	tree.insert(tree.root, data)
}
func (tree *BTree) insert(curNode *BTreeNode, data int) (overflowKey int, leftChild *BTreeNode, rightChild *BTreeNode) {
	if curNode.childs[0] != nil {
		for i := 0; i <= curNode.keyLen; i++ {
			if data < curNode.keys[i] || i == curNode.keyLen {
				overflowKey, leftChild, rightChild = tree.insert(curNode.childs[i], data)
				if leftChild == nil {
					return
				}
				//有上溢出，则将上溢出的节点插入到当前节点
				insertPosition := curNode.insertKey(overflowKey)
				curNode.childs[insertPosition] = leftChild
				//将右孩子插入子子节点中
				for i := curNode.keyLen - 1; i > insertPosition; i-- {
					curNode.childs[i+1] = curNode.childs[i]
				}
				curNode.childs[insertPosition+1] = rightChild
				for i := 0; i <= curNode.keyLen; i++ {
					curNode.childs[i].parent = curNode
				}
				//如果插入之后，没有溢出则返回
				if curNode.keyLen < tree.rank {
					return 0, nil, nil
				}
				//处理溢出
				overflowKeyIdx := tree.rank / 2
				overflowKey = curNode.keys[overflowKeyIdx]
				leftChild = &BTreeNode{keys: make([]int, tree.rank), keyLen: overflowKeyIdx, childs: make([]*BTreeNode, tree.rank+1)}
				//赋值左孩子的值
				copy(leftChild.keys, curNode.keys[:overflowKeyIdx])
				//赋值左孩子的子节点
				copy(leftChild.childs, curNode.childs[:overflowKeyIdx+1])
				for i := 0; i <= overflowKeyIdx; i++ {
					leftChild.childs[i].parent = leftChild
				}
				//复制右孩子
				rightChild = &BTreeNode{keys: make([]int, tree.rank), keyLen: tree.rank - (overflowKeyIdx + 1), childs: make([]*BTreeNode, tree.rank+1)}
				//复制右孩子的key
				copy(rightChild.keys, curNode.keys[overflowKeyIdx+1:])
				//复制右孩子的值
				copy(rightChild.childs, curNode.childs[overflowKeyIdx+1:])
				for i := 0; i < (tree.rank - overflowKeyIdx); i++ {
					rightChild.childs[i].parent = rightChild
				}
				if curNode == tree.root {
					tree.root = &BTreeNode{keys: make([]int, tree.rank), childs: make([]*BTreeNode, tree.rank+1)}
					tree.root.keys[0] = overflowKey
					tree.root.childs[0] = leftChild
					tree.root.childs[1] = rightChild
					leftChild.parent = tree.root
					rightChild.parent = tree.root
					tree.root.keyLen = 1
					return 0, nil, nil
				}
				return overflowKey, leftChild, rightChild
			}
		}
	}

	curNode.insertKey(data)
	if curNode.keyLen < tree.rank {
		return
	}
	overflowKeyIdx := tree.rank / 2
	overflowKey = curNode.keys[overflowKeyIdx]
	leftChild = &BTreeNode{keys: make([]int, tree.rank), keyLen: overflowKeyIdx, childs: make([]*BTreeNode, tree.rank+1)}
	copy(leftChild.keys, curNode.keys[:overflowKeyIdx])

	rightChild = &BTreeNode{keys: make([]int, tree.rank), keyLen: tree.rank - (overflowKeyIdx + 1), childs: make([]*BTreeNode, tree.rank+1)}
	copy(rightChild.keys, curNode.keys[overflowKeyIdx+1:])
	if curNode == tree.root {
		tree.root = &BTreeNode{keys: make([]int, tree.rank), childs: make([]*BTreeNode, tree.rank+1)}
		tree.root.keys[0] = overflowKey
		tree.root.childs[0] = leftChild
		tree.root.childs[1] = rightChild
		leftChild.parent = tree.root
		rightChild.parent = tree.root
		tree.root.keyLen = 1
		return 0, nil, nil
	}
	return
}
func (node *BTreeNode) insertKey(insertKey int) (insertPosition int) {
	//insert sort
	if node.keyLen == 0 {
		node.keys[0] = insertKey
		node.keyLen = 1
		return 0
	}
	var i = 0
	for i = node.keyLen; i > 0; i-- {
		if node.keys[i-1] > insertKey {
			node.keys[i] = node.keys[i-1]
			continue
		}
		break
	}
	node.keys[i] = insertKey
	node.keyLen++
	return i
}
func (tree *BTree) Delete(data int) bool {
	return tree.delete(tree.root, data)
}
func (tree *BTree) delete(node *BTreeNode, data int) bool {
	if node == nil {
		return false
	}
	pos, isEqual := tree.find(node, data)
	if isEqual {
		//左边的右孩子
		//如果找到的节点没有孩子，则直接删除，然后调整
		if node.childs[pos] == nil {
			tree.shift(node, pos)
			return true
		}
		//如果有左孩子，没有右孩子，则替换左孩子节点
		if node.childs[pos].childs[0] == nil {
			node.keys[pos] = node.childs[pos].keys[node.childs[pos].keyLen-1]
			tree.shift(node.childs[pos], node.childs[pos].keyLen-1)
			return true
		}
		//找到左孩子的最右的节点
		replaceNode := tree.replace(node.childs[pos].childs[node.childs[pos].keyLen])
		node.keys[pos] = replaceNode.keys[replaceNode.keyLen-1]
		tree.shift(replaceNode, replaceNode.keyLen-1)
		return true
	}
	//往左搜索
	if pos == -1 {
		return tree.delete(node.childs[0], data)
	}
	//往右搜索
	return tree.delete(node.childs[pos+1], data)
}
func (tree *BTree) find(node *BTreeNode, data int) (int, bool) {
	low := 0
	high := node.keyLen - 1
	mid := 0
	for low <= high {
		mid = (low + high) / 2
		if node.keys[mid] == data {
			return mid, true
		}
		if node.keys[mid] > data {
			high = mid - 1
			continue
		}
		low = mid + 1
	}
	//2 4 6 8 10
	if low > high {
		return low - 1, false
	}
	return high, false
}
func (tree *BTree) replace(node *BTreeNode) *BTreeNode {
	//找到左边的右孩子
	if node.childs[0] == nil {
		return node
	}
	return tree.replace(node.childs[node.keyLen])
}
func (tree *BTree) shift(node *BTreeNode, pos int) {
	if node == nil {
		return
	}
	replaceKey := node.keys[pos]
	for i := pos; i < node.keyLen; i++ {
		node.keys[i] = node.keys[i+1]
	}
	node.keyLen -= 1
	//如果调整到root节点
	if node.keyLen == 0 && node == tree.root {
		tree.root = node.childs[0]
		return
	}
	//如果调整的node是root，则不必判断是否满足临界条件
	if node == tree.root {
		return
	}
	crisis := int(math.Ceil(float64(tree.rank)/2 - 1))
	//判断是否小于临界节点数
	if node.keyLen < crisis {
		//向左节点借
		//找到第一个父节点小于replaceKey的位置，该位置对应的child为node的左兄弟
		ppos, equal := tree.find(node.parent, replaceKey)
		if equal {
			ppos -= 1
		}
		var leftChild *BTreeNode
		var rightChild *BTreeNode
		var browerNode *BTreeNode
		direction := 0
		//该节点在父节点的左边，所以只能去右节点去借
		if ppos < 0 {
			//如果借的节点，借完之后小于临界，则左，父亲，右节点进行合并
			leftChild = node
			rightChild = node.parent.childs[1]
			browerNode = rightChild
			ppos = 0
			direction = 0
		} else {
			leftChild = node.parent.childs[ppos]
			rightChild = node
			browerNode = leftChild
			direction = 1
		}
		//节点合并
		//如果借的节点，借完之后小于临界，则左，父亲，右节点进行合并
		if browerNode.keyLen-1 < crisis {
			newNode := &BTreeNode{keys: make([]int, tree.rank), childs: make([]*BTreeNode, tree.rank+1)}
			newNode.parent = leftChild.parent
			copy(newNode.keys, leftChild.keys[:leftChild.keyLen])
			newNode.keyLen += leftChild.keyLen
			copy(newNode.childs, leftChild.childs[:leftChild.keyLen+1])
			newNode.keys[newNode.keyLen] = newNode.parent.keys[ppos]
			newNode.keyLen += 1
			copy(newNode.keys[newNode.keyLen:], rightChild.keys[:rightChild.keyLen])
			copy(newNode.childs[newNode.keyLen:], rightChild.childs[:rightChild.keyLen+1])
			newNode.keyLen += rightChild.keyLen

			//将新的节点挂载到父节点,父节点缩减一个孩子，因为父节点有一个key并入了
			for i := ppos; i < newNode.parent.keyLen; i++ {
				newNode.parent.childs[i] = newNode.parent.childs[i+1]
			}
			newNode.parent.childs[ppos] = newNode

			//修改newNode childs的父亲节点
			if newNode.childs[0] != nil && newNode.keyLen > 0 {
				for i := 0; i <= newNode.keyLen; i++ {
					newNode.childs[i].parent = newNode
				}
			}
			//继续调整父节点
			tree.shift(newNode.parent, ppos)
			return
		}
		if direction == 0 {
			//能从右边借到节点
			leftChild.keys[leftChild.keyLen] = leftChild.parent.keys[0]
			leftChild.keyLen += 1
			leftChild.parent.keys[0] = rightChild.keys[0]
			//分离出右孩子的孩子们的最左边的孩子,挂载到左孩子的孩子们的最右边
			splitChild := rightChild.childs[0]
			for i := 0; i < rightChild.keyLen-1; i++ {
				rightChild.keys[i] = rightChild.keys[i+1]
			}
			//调整右孩子的孩子们的数量向左平移一位
			for i:=0;i<rightChild.keyLen;i++{
				rightChild.childs[i]=rightChild.childs[i+1]
			}
			rightChild.keyLen -= 1
			leftChild.childs[leftChild.keyLen] = splitChild
			if splitChild != nil {
				splitChild.parent = leftChild
			}
			return
		}
		if direction == 1 {
			//能从左边借到节点
			//将右孩子的keys一起往后移动一位，留出第一个位置
			for i := rightChild.keyLen; i > 0; i-- {
				rightChild.keys[i] = rightChild.keys[i-1]
			}
			rightChild.keys[0] = rightChild.parent.keys[ppos]
			rightChild.keyLen += 1
			//将左边节点的最右边key给当前的parent
			leftChild.parent.keys[ppos] = leftChild.keys[leftChild.keyLen-1]
			//将左孩子的孩子们的右孩子，挂载到右孩子的最左边，因为key调整
			splitChild := leftChild.childs[leftChild.keyLen]
			leftChild.childs[leftChild.keyLen] = nil
			leftChild.keyLen -= 1
			//将splitChild挂载到右孩子的孩子们的最左边
			for i := rightChild.keyLen; i > 0; i-- {
				rightChild.childs[i] = rightChild.childs[i-1]
			}
			rightChild.childs[0] = splitChild
			if splitChild != nil {
				splitChild.parent = rightChild
			}
		}

	}
}
