package bst

import "fmt"

const (
	POSITION_INTT   = 0
	POSITION_LEFT   = 1
	POSITION_RIGHT  = 2
	DIRECTION_LEFT  = 1
	DIRECTION_RIGHT = 2
)

type Bst struct {
	root *TreeNode
}
type TreeNode struct {
	data  int
	left  *TreeNode
	right *TreeNode
}

func (bst *Bst) Insert(data int) {
	if bst.root == nil {
		bst.root = &TreeNode{data: data}
		return
	}
	bst.insert(bst.root, data)
}
func (bst *Bst) insert(node *TreeNode, data int) *TreeNode {
	if node == nil {
		return &TreeNode{data: data}
	}
	if data < node.data {
		node.left = bst.insert(node.left, data)
	} else {
		node.right = bst.insert(node.right, data)
	}
	return node
}
func (bst *Bst) Delete(data int) bool {
	if bst.root == nil {
		return false
	}
	return bst.delete(nil, bst.root, POSITION_INTT, data)
}
func (bst *Bst) delete(nodeParent *TreeNode, curNode *TreeNode, position int, data int) bool {
	if curNode == nil {
		return false
	}
	if curNode.data == data {
		if curNode.left != nil {
			replacedNodeParent, replacedNode := bst.replaceOnDelete(curNode.left, curNode.left.right, DIRECTION_RIGHT)
			// not exist nil,nil
			//not found replaced node on direction right, so curNode change with curNode.left node
			if replacedNode == nil {
				curNode.data = curNode.left.data
				curNode.left = curNode.left.left
			} else {
				curNode.data = replacedNode.data
				replacedNodeParent.right = replacedNode.left
			}
			return true
		}
		if curNode.right != nil {
			replacedNodeParent, replacedNode := bst.replaceOnDelete(curNode.right, curNode.right.left, DIRECTION_LEFT)
			// not exist nil,nil
			//not found replaced node on direction right, so curNode change with curNode.left node
			if replacedNode == nil {
				curNode.data = curNode.right.data
				curNode.right = curNode.right.right
			} else {
				curNode.data = replacedNode.data
				replacedNodeParent.left = replacedNode.right
			}
			return true
		}
		// case curNode.right==nil && curNode.left==nil
		if position == POSITION_INTT {
			//represent curNode is root
			bst.root = nil
		}
		if position == POSITION_LEFT {
			nodeParent.left = nil
		}
		if position == POSITION_RIGHT {
			nodeParent.right = nil
		}
		return true
	}
	flag := false
	flag = bst.delete(curNode, curNode.left, POSITION_LEFT, data)
	if flag {
		return true
	}
	return bst.delete(curNode, curNode.right, POSITION_RIGHT, data)
}
func (bst *Bst) replaceOnDelete(parentNode *TreeNode, curNode *TreeNode, direction int) (replacedNodeParent *TreeNode, replacedNode *TreeNode) {
	replacedNodeParent = parentNode
	if curNode == nil {
		return replacedNodeParent, nil
	}
	if curNode.left == nil && curNode.right == nil {
		return replacedNodeParent, curNode
	}
	if direction == DIRECTION_LEFT {
		if curNode.left == nil {
			return replacedNodeParent, curNode
		}
		return bst.replaceOnDelete(curNode, curNode.left, direction)
	}
	if direction == DIRECTION_RIGHT {
		if curNode.right == nil {
			return replacedNodeParent, curNode
		}
		return bst.replaceOnDelete(curNode, curNode.left, direction)
	}
	return nil, nil
}
func LDR(node *TreeNode) {
	if node == nil {
		return
	}
	LDR(node.left)
	fmt.Print(node.data, " ")
	LDR(node.right)
}
