package redtree

const (
	red            = 1
	black          = 2
	PositionInit   = 0
	PositionLeft   = 1
	PositionRight  = 2
	DirectionLeft  = 1
	DirectionRight = 2
)

type RedBlackTree struct {
	root *TreeNode
}
type TreeNode struct {
	color  int
	data   int
	left   *TreeNode
	right  *TreeNode
	parent *TreeNode
}

func (tree *RedBlackTree) Insert(data int) {
	if tree.root == nil {
		tree.root = &TreeNode{color: black, data: data}
	} else {
		tree.insert(tree.root, data)
	}
}
func (tree *RedBlackTree) insert(curNode *TreeNode, data int) {
	if data <= curNode.data {
		if curNode.left == nil {
			curNode.left = &TreeNode{data: data, color: red, parent: curNode}
			tree.shift(curNode.left)
			return
		} else {
			tree.insert(curNode.left, data)
			tree.shift(curNode)
			return
		}

	} else {
		if curNode.right == nil {
			curNode.right = &TreeNode{data: data, color: red, parent: curNode}
			tree.shift(curNode.right)
			return
		} else {
			tree.insert(curNode.right, data)
			tree.shift(curNode)
			return
		}
	}
}
func (tree *RedBlackTree) shift(child *TreeNode) {
	if child.parent == nil {
		child.color = black
		return
	}
	if child.color == red && child.parent != nil && child.parent.color == red {
		//if child's parent color is red, then draw parent's color as red
		if child.parent.parent == nil {
			child.parent.color = black
		} else {

			childPosition := PositionInit
			childParentPosition := PositionInit
			childParentParentPosition := PositionInit
			var uncleNode *TreeNode
			if child.data <= child.parent.data {
				childPosition = PositionLeft
			} else {
				childPosition = PositionRight
			}

			if child.parent.data <= child.parent.parent.data {
				childParentPosition = PositionLeft
				uncleNode = child.parent.parent.right
			} else {
				childParentPosition = PositionRight
				uncleNode = child.parent.parent.left
			}
			A := child
			B := child.parent
			C := child.parent.parent
			D := child.parent.parent.parent
			if D != nil {
				if C.data <= D.data {
					childParentParentPosition = PositionLeft
				} else {
					childParentParentPosition = PositionRight
				}
			}

			// if uncle's color is black, then rotate
			if uncleNode == nil || uncleNode.color == black {

				//LL
				/*
					    C
					  B
					A
				*/
				if childPosition == PositionLeft && childParentPosition == PositionLeft {
					C.left = B.right
					if C.left != nil {
						C.left.parent = C
					}
					B.right = C
					C.parent = B
					if D == nil {
						tree.root = B
						B.parent = nil
					} else {
						if childParentParentPosition == PositionLeft {
							D.left = B
						} else {
							D.right = B
						}
						B.parent = D
					}
					B.color = black
					B.left.color = red
					B.right.color = red
				}
				//LR
				/*
						C
					B
						A

				*/
				if childPosition == PositionRight && childParentPosition == PositionLeft {
					B.right = A.left
					if B.right != nil {
						B.right.parent = B
					}
					A.left = B
					B.parent = A
					C.left = A.right
					if C.left != nil {
						C.left.parent = C
					}
					A.right = C
					C.parent = A
					B.right = nil
					C.left = nil
					if D == nil {
						tree.root = A
						A.parent = nil
					} else {
						if childParentParentPosition == PositionLeft {
							D.left = A
						} else {
							D.right = A
						}
						A.parent = D
					}
					A.color = black
					A.left.color = red
					A.right.color = red
				}
				//RR
				/*
					C
						B
							A
				*/
				if childPosition == PositionRight && childParentPosition == PositionRight {
					C.right = B.left
					if C.right != nil {
						C.right.parent = C
					}
					B.left = C
					C.parent = B
					if D == nil {
						tree.root = B
						B.parent = nil
					} else {
						if childParentParentPosition == PositionLeft {
							D.left = B
						} else {
							D.right = B
						}
						B.parent = D
					}
					B.color = black
					B.left.color = red
					B.right.color = red
				}
				//RL
				/*
						C
							B
					    A
				*/
				if childPosition == PositionLeft && childParentPosition == PositionRight {
					C.right = A.left
					if C.right != nil {
						C.right.parent = C
					}
					A.left = C
					C.parent = A
					B.left = A.right
					if B.left != nil {
						B.left.parent = B
					}
					A.right = B
					B.parent = A
					if D == nil {
						tree.root = A
						A.parent = nil
					} else {
						if childParentParentPosition == PositionLeft {
							D.left = A
						} else {
							D.right = A
						}
						A.parent = D
					}
					A.color = black
					A.left.color = red
					A.right.color = red
				}
			} else {
				//if uncle 's color is red
				uncleNode.color = black
				B.color = black
				C.color = red
			}
		}
	}
	return
}
func (tree *RedBlackTree) Delete(data int) bool {
	return tree.delete(tree.root, data)
}
func (tree *RedBlackTree) delete(curNode *TreeNode, data int) bool {
	if curNode == nil {
		return false
	}
	if curNode.data == data {
		//由于root的color一直是黑色的，所以当curNode为红色是一定不是黑色节点
		//情况1：左右树有一个为空，则不为空的子节点肯定为1个或者0个，如果被删除节点是红色，则其左右孩子节点肯定是都是
		//空节点，如果被删除节点是黑色，那么不为空的的节点的颜色肯定是红色，则用该节点直接替换，并修改成被删除的节点的颜色即可
		if curNode.left == nil || curNode.right == nil {
			//删除curNode节点
			//找到curNode在父节点的位置
			if curNode.left == nil && curNode.right == nil {
				if tree.root == curNode {
					tree.root = nil
					return true
				}
				if curNode.data <= curNode.parent.data {
					curNode.parent.left = nil
				} else {
					curNode.parent.right = nil
				}
				if curNode.color == black {
					tree.doubleBlackShift(curNode)
					return true
				}
			} else if curNode.left == nil {
				if curNode == tree.root {
					tree.root = curNode.right
					tree.root.parent = nil
					tree.root.color = black
					return true
				}
				if curNode.data <= curNode.parent.data {
					curNode.parent.left = curNode.right
					curNode.parent.left.color = curNode.color
					curNode.parent.left.parent = curNode.parent

				} else {
					curNode.parent.right = curNode.right
					curNode.parent.right.color = curNode.color
					curNode.parent.right.parent = curNode.parent
				}
			} else if curNode.right == nil {
				if curNode == tree.root {
					tree.root = curNode.left
					tree.root.parent = nil
					tree.root.color = black
					return true
				}
				if curNode.data <= curNode.parent.data {
					curNode.parent.left = curNode.left
					curNode.parent.left.color = curNode.color
					curNode.parent.left.parent = curNode.parent
				} else {
					curNode.parent.right = curNode.left
					curNode.parent.right.color = curNode.color
					curNode.parent.right.parent = curNode.parent
				}
			}
		}
		//如果curNode的color是红色，并且左右节点多不为空，则需要遍历左子树的最右孩子或者右子树的最左孩子
		//进行替换，并根据被替换的节点的颜色进行调制。
		if curNode.left != nil && curNode.right != nil {
			foundNode := tree.findReplacedNode(curNode.left, DirectionRight)
			if foundNode == nil {
				//如果没有找到，则curNode.left节点替换CurNode节点
				if curNode.data <= curNode.parent.data {
					curNode.parent.left = curNode.left
					curNode.parent.left.parent = curNode.parent
				} else {
					curNode.parent.right = curNode.left
					curNode.parent.right.parent = curNode.parent
				}
			} else {
				curNode.data = foundNode.data
				//直接删除，并把挂载foundData的左节点
				foundNode.parent.right = foundNode.left

				if foundNode.color == black {
					tree.doubleBlackShift(foundNode)
				}

			}
		}
		return true
	} else if data < curNode.data {
		return tree.delete(curNode.left, data)
	} else {
		return tree.delete(curNode.right, data)
	}
}
func (tree *RedBlackTree) doubleBlackShift(curNode *TreeNode) {
	if curNode.parent == nil {
		//如果调整
		return
	}
	//出现双黑节点，则调整
	var brotherNode *TreeNode
	brotherDirection := PositionInit
	if curNode.data <= curNode.parent.data {
		brotherNode = curNode.parent.right
		brotherDirection = PositionRight
	} else {
		brotherNode = curNode.parent.left
		brotherDirection = PositionLeft
	}
	//情况1：如果兄弟是黑色，且左右孩子有一个红色,foundData的兄弟节点肯定不会为空，因为删除的foundData的颜色为黑色
	if brotherNode.color == black {

		if brotherDirection == DirectionLeft && brotherNode.left != nil && brotherNode.left.color == red {
			//LL型，调整成支架型，中间节点随自己的父节点，支架的左右节点调整为黑色
			//    C
			//  B
			//A
			B := brotherNode
			C := brotherNode.parent
			A := brotherNode.left
			D := C.parent
			B.color = C.color
			brotherNodeParentParent := C.parent
			C.left = B.right
			C.left.parent = C
			B.right = C
			C.parent = B
			B.left.color = black
			B.right.color = black
			if brotherNodeParentParent == nil {
				tree.root = B
				B.parent = nil
			} else {
				if B.data <= D.data {
					D.left = B
					B.parent = D
				} else {
					D.right = B
					B.parent = D
				}
			}
			A.color = red
			C.color = red
		} else if brotherDirection == DirectionLeft && brotherNode.right != nil && brotherNode.right.color == red {
			//LR型，调整为支架型
			//   C
			// B
			//   A
			A := brotherNode.right
			B := brotherNode
			C := brotherNode.parent
			D := brotherNode.parent.parent
			C.left = A.right
			A.right.parent = C
			B.right = A.left
			B.right.parent = B
			A.left = B
			B.parent = A
			A.right = C
			C.parent = A
			A.color = C.color
			B.color = black
			C.color = black
			if D == nil {
				tree.root = A
				A.parent = nil
			} else {
				if A.data <= D.data {
					D.left = A
					A.parent = D
				} else {
					D.right = A
					A.parent = D
				}
			}
		} else if brotherDirection == DirectionRight && brotherNode.right != nil && brotherNode.right.color == red {
			//RR 调整成支架型，中间节点随自己的父节点，支架的左右节点调整为黑色
			//C
			//  B
			//    A
			B := brotherNode
			C := brotherNode.parent
			//A := brotherNode.right
			D := C.parent
			B.color = C.color
			brotherNodeParentParent := C.parent
			C.right = B.left
			if C.right != nil {
				C.right.parent = C
			}
			B.left = C
			C.parent = B
			B.left.color = black
			B.right.color = black
			if brotherNodeParentParent == nil {
				tree.root = brotherNode
				B.parent = nil
			} else {
				if B.data <= D.data {
					D.left = B
					B.parent = D
				} else {
					D.right = B
					B.parent = D
				}
			}

		} else if brotherDirection == DirectionRight && brotherNode.left != nil && brotherNode.left.color == red {
			//RL型，调整为支架型
			//   C
			// 		B
			//   A
			A := brotherNode.left
			B := brotherNode
			C := brotherNode.parent
			D := brotherNode.parent.parent
			C.right = A.left
			if C.right != nil {
				C.right.parent = C
			}
			B.left = A.right
			if B.left!=nil{
				B.left.parent = B
			}
			A.right = B
			B.parent = A
			A.left = C
			C.parent = A
			A.color = C.color
			B.color = black
			C.color = black
			if D == nil {
				tree.root = A
				A.parent = nil
			} else {
				if A.data <= D.data {
					D.left = A
					A.parent = D
				} else {
					D.right = A
					A.parent = D
				}
			}
		} else {
			//如果brother为黑色，左右孩子都是黑色，则换色
			//如果brother的父节点为红色，则换色后不调整
			originParentColor := brotherNode.parent.color
			brotherNode.color = red
			brotherNode.parent.color = black
			if originParentColor == black {
				//brother的父节点为双黑节点
				tree.doubleBlackShift(brotherNode.parent)
			}
		}
	} else {
		if brotherDirection == DirectionLeft {
			//如果brother节点为红节点，则mirror后，再进行调整
			//LL 型
			//  C
			// B
			//A
			B := brotherNode
			//A := B.left
			C := B.parent
			D := C.parent
			C.left = B.right
			if C.left != nil {
				C.left.parent = C
			}
			B.right = C
			C.parent = B
			if D == nil {
				tree.root = B
				B.parent = nil
			} else {
				if B.data <= D.data {
					D.left = B
					B.parent = D
				} else {
					D.right = B
					B.parent = D
				}
			}
			B.color = C.color
			C.color = red
			//标记继续调整
			tree.doubleBlackShift(curNode)
		} else {
			//如果brother节点为红节点，则mirror后，再进行调整
			//RR 型
			//  C
			// 	  B
			//		A
			B := brotherNode
			//A := B.left
			C := B.parent
			D := C.parent
			C.right = B.left
			if C.right != nil {
				C.right.parent = C
			}
			B.left = C
			C.parent = B
			if D == nil {
				tree.root = B
				B.parent = nil
			} else {
				if B.data <= D.data {
					D.left = B
					B.parent = D
				} else {
					D.right = B
					B.parent = D
				}
			}
			B.color = C.color
			C.color = red
			//标记继续调整
			tree.doubleBlackShift(curNode)
		}
	}
}
func (tree *RedBlackTree) findReplacedNode(curNode *TreeNode, direction int) *TreeNode {
	if direction == DirectionLeft {
		if curNode.left == nil {
			return curNode
		} else {
			return tree.findReplacedNode(curNode.left, direction)
		}
	} else {
		if curNode.right == nil {
			return curNode
		} else {
			return tree.findReplacedNode(curNode.right, direction)
		}
	}
}
