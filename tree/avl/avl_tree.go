package avl

import (
	"fmt"
)

const (
	PositionInit  = 0
	PositionLeft  = 1
	PositionRight = 2
)

type AVLNode struct {
	left      *AVLNode
	right     *AVLNode
	data      int
	leftHight int
	righHight int
}
type AVLTree struct {
	root *AVLNode
}

func NewAvlTree() *AVLTree {
	return &AVLTree{}
}
func (tree *AVLTree) Put(data int) {
	if tree.root == nil {
		tree.root = &AVLNode{data: data, leftHight: 0, righHight: 0}
		return
	}
	tree.insert(nil, PositionInit, tree.root, data)
}
func (tree *AVLTree) del(curRootParent *AVLNode, curRootPosition int, curRoot *AVLNode, data int) (isChange bool) {
	if curRoot == nil {
		return false
	}
	if curRoot.data == data {
		//如果curRoot是叶子节点
		if curRoot.left == nil && curRoot.right == nil {
			if curRootParent == nil {
				tree.root = nil
				return false
			}
			if curRootPosition == PositionLeft {
				curRootParent.left = nil
				curRootParent.leftHight -= 1
			} else {
				curRootParent.right = nil
				curRootParent.righHight -= 1
			}
		} else if curRoot.leftHight > curRoot.righHight {
			isReplaced := tree.replaceRight(curRoot, curRoot.left, curRoot.left.right)
			if !isReplaced {
				curRoot.data = curRoot.left.data
				curRoot.left = nil
				curRoot.leftHight = 0
			}
		} else {
			isReplaced := tree.replaceRight(curRoot, curRoot.right, curRoot.right.left)
			if !isReplaced {
				curRoot.data = curRoot.right.data
				curRoot.right = nil
				curRoot.righHight = 0
			}

		}
		return true
	} else {
		if data < curRoot.data {
			isChange = tree.del(curRoot, PositionLeft, curRoot.left, data)
			if curRoot.left != nil {
				curRoot.leftHight = max(curRoot.left.leftHight, curRoot.left.righHight) + 1
			}

			if !isChange {
				return isChange
			}
			if abs(curRoot.leftHight, curRoot.righHight) >= 2 {
				if curRoot.right.righHight > curRoot.right.leftHight {
					// A
					//    B
					//       C
					A := curRoot
					B := curRoot.right
					//C := curRoot.right.right

					//step1:将B节点作为中心节点，放到A的位置
					if curRootParent == nil {
						tree.root = B
					} else {
						if curRootPosition == PositionLeft {
							curRootParent.left = B
						} else {
							curRootParent.right = B
						}
					}
					//step2:将A的右孩子作为B的左孩子，B的左孩子作为A的右孩子
					A.right = B.left
					B.left = A
					//step3:调整各自的leftHeight和rightHeight
					if A.right != nil {
						A.righHight = max(A.right.leftHight, A.right.righHight) + 1
					} else {
						A.righHight = 0
					}

					B.leftHight = max(B.left.leftHight, B.left.righHight) + 1
					if curRootParent != nil {
						if curRootPosition == PositionLeft {
							curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
						} else {
							curRootParent.righHight = max(curRootParent.right.leftHight, curRootParent.right.righHight) + 1
						}
					}
					return false
				} else {
					//RL

					//  A
					//     B
					//  C
					A := curRoot
					B := curRoot.right
					C := curRoot.right.left
					//step1:将C放到A的位置
					if curRootParent == nil {
						tree.root = C
					} else {
						if curRootPosition == PositionLeft {
							curRootParent.left = C
						} else {
							curRootParent.right = C
						}
					}
					//step2:将A的right指向C的left
					A.right = C.left
					//step3:将B的left指向C的right
					B.left = C.right
					//step4:将A作为C的左孩子，B作为C的右孩子
					C.left = A
					C.right = B
					//step5:调整相关节点的高度
					if A.right != nil {
						A.righHight = max(A.right.leftHight, A.right.righHight) + 1
					} else {
						A.righHight = 0
					}
					if B.left != nil {
						B.leftHight = max(B.left.leftHight, B.left.righHight) + 1
					} else {
						B.leftHight = 0
					}
					C.leftHight = max(C.left.leftHight, C.left.righHight) + 1
					C.righHight = max(C.right.leftHight, C.right.righHight) + 1
					if curRootParent != nil {
						if curRootPosition == PositionLeft {
							curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
						} else {
							curRootParent.righHight = max(curRootParent.leftHight, curRootParent.righHight) + 1
						}
					}
				}
				return false
			}
			return isChange
		} else {
			isChange = tree.del(curRoot, PositionRight, curRoot.right, data)
			if curRoot.right != nil {
				curRoot.righHight = max(curRoot.right.leftHight, curRoot.right.righHight) + 1
			}
			if !isChange {
				return isChange
			}
			//curRoot.righHight += 1
			//check and shift
			if abs(curRoot.leftHight, curRoot.righHight) >= 2 {
				if curRoot.left.leftHight > curRoot.left.righHight {
					//LL
					//      A
					//    B
					//  C
					A := curRoot
					B := curRoot.left
					//	C:=curRoot.left.left

					//setp1:将当前curRoot摘除，将curRoot的左孩子挂到curRoot当前的位置
					if curRootParent == nil {
						tree.root = B
					} else {
						//判断curRoot在左孩子
						if curRootPosition == PositionLeft {
							curRootParent.left = B
						} else {
							curRootParent.right = B
						}
					}
					//setp2:将curRoot左孩子节点的右孩子挂到curRoot的左孩子，此时curRoot的left节点和curRoot脱钩
					A.left = B.right
					//setp3:将curRoot挂到，先前的left节点的右节点
					B.right = A

					//step3:调整各自的leftHeight和rightHeight
					if A.left != nil {
						A.leftHight = max(A.left.leftHight, A.left.righHight) + 1
					} else {
						A.leftHight = 0
					}
					B.righHight = max(B.right.leftHight, B.right.righHight) + 1
					if curRootParent != nil {
						if curRootPosition == PositionLeft {
							curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
						} else {
							curRootParent.righHight = max(curRootParent.right.leftHight, curRootParent.right.righHight) + 1
						}
					}
					return false
				} else {
					//LR
					//    A
					//  B
					//    C
					//setp1:将C节点挂到A的位置
					A := curRoot
					B := curRoot.left
					C := curRoot.left.right
					if curRootParent == nil {
						//如果A节点是root节点
						tree.root = C
					} else {
						//判断curRoot在左孩子
						if curRootPosition == PositionLeft {
							curRootParent.left = C
						} else {
							curRootParent.right = C
						}
					}
					//step2:将C的左节点挂到B的右孩子
					B.right = C.left
					//setp3:将C的右节点挂到A的左孩子
					A.left = C.right
					//step4:将B节点作为C的左孩子
					C.left = B
					//step5:将A节点作为C的右孩子
					C.right = A
					//step6:调整相关节点的高度
					//if B != nil {
					if B.right != nil {
						B.righHight = max(B.right.leftHight, B.right.righHight) + 1
					} else {
						B.righHight = 0
					}
					//}
					//if A != nil {
					if A.left != nil {
						A.leftHight = max(A.left.leftHight, A.left.righHight) + 1
					} else {
						A.leftHight = 0
					}
					//}
					C.leftHight = max(C.left.leftHight, C.left.righHight) + 1
					C.righHight = max(C.right.leftHight, C.right.righHight) + 1
					if curRootParent != nil {
						if curRootPosition == PositionLeft {
							curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
						} else {
							curRootParent.righHight = max(curRootParent.right.leftHight, curRootParent.right.righHight) + 1
						}
					}
				}
				return false
			}

		}
		return isChange
	}

}
func (tree *AVLTree) replaceRight(needReplaceNode *AVLNode, curRootParent *AVLNode, curRoot *AVLNode) bool {
	if curRoot == nil {
		return false
	}

	if curRoot.left != nil {
		replaced := tree.replaceRight(needReplaceNode, curRoot, curRoot.left)
		//if replaced {
		//	curRootParent.righHight = max(curRootParent.right.leftHight, curRootParent.right.righHight)
		//}
		return replaced
	} else {
		needReplaceNode.data = curRoot.data
		if curRootParent.left.right != nil {
			curRootParent.left = curRootParent.left.right
			curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
		} else {
			curRootParent.left = nil
			curRootParent.leftHight = 0
		}

		return true
	}
}
func (tree *AVLTree) replaceLeft(needReplaceNode *AVLNode, curRootParent *AVLNode, curRoot *AVLNode) bool {
	if curRoot == nil {
		return false
	}

	if curRoot.right != nil {
		replaced := tree.replaceLeft(needReplaceNode, curRoot, curRoot.right)
		//if replaced {
		//	curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight)
		//}
		return replaced
	} else {
		needReplaceNode.data = curRoot.data
		if curRootParent.right.left != nil {
			curRootParent.right = curRootParent.right.left
			curRootParent.righHight = max(curRootParent.right.leftHight, curRootParent.right.righHight) + 1
		} else {
			curRootParent.right = nil
			curRootParent.righHight = 0
		}

		return true
	}
}
func (tree *AVLTree) insert(curRootParent *AVLNode, curRootPosition int, curRoot *AVLNode, data int) (isChange bool) {
	if data <= curRoot.data {

		if curRoot.left != nil {
			isChange = tree.insert(curRoot, PositionLeft, curRoot.left, data)
			curRoot.leftHight = max(curRoot.left.leftHight, curRoot.left.righHight) + 1
			if !isChange {
				return isChange
			}

			//curRoot.leftHight += 1
			//check and shift
			if abs(curRoot.leftHight, curRoot.righHight) >= 2 {
				insertPosition := PositionInit
				if data <= curRoot.left.data {
					insertPosition = PositionLeft
				} else {
					insertPosition = PositionRight
				}
				if insertPosition == PositionLeft {
					//LL
					//      A
					//    B
					//  C
					A := curRoot
					B := curRoot.left
					//	C:=curRoot.left.left

					//setp1:将当前curRoot摘除，将curRoot的左孩子挂到curRoot当前的位置
					if curRootParent == nil {
						tree.root = B
					} else {
						//判断curRoot在左孩子
						if curRootPosition == PositionLeft {
							curRootParent.left = B
						} else {
							curRootParent.right = B
						}
					}
					//setp2:将curRoot左孩子节点的右孩子挂到curRoot的左孩子，此时curRoot的left节点和curRoot脱钩
					A.left = B.right
					//setp3:将curRoot挂到，先前的left节点的右节点
					B.right = A

					//step3:调整各自的leftHeight和rightHeight
					if A.left != nil {
						A.leftHight = max(A.left.leftHight, A.left.righHight) + 1
					} else {
						A.leftHight = 0
					}
					B.righHight = max(B.right.leftHight, B.right.righHight) + 1
					if curRootParent != nil {
						if curRootPosition == PositionLeft {
							curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
						} else {
							curRootParent.righHight = max(curRootParent.right.leftHight, curRootParent.right.righHight) + 1
						}
					}
				} else {
					//LR
					//    A
					//  B
					//    C
					//setp1:将C节点挂到A的位置
					A := curRoot
					B := curRoot.left
					C := curRoot.left.right
					if curRootParent == nil {
						//如果A节点是root节点
						tree.root = C
					} else {
						//判断curRoot在左孩子
						if curRootPosition == PositionLeft {
							curRootParent.left = C
						} else {
							curRootParent.right = C
						}
					}
					//step2:将C的左节点挂到B的右孩子
					B.right = C.left
					//setp3:将C的右节点挂到A的左孩子
					A.left = C.right
					//step4:将B节点作为C的左孩子
					C.left = B
					//step5:将A节点作为C的右孩子
					C.right = A
					//step6:调整相关节点的高度
					//if B != nil {
					if B.right != nil {
						B.righHight = max(B.right.leftHight, B.right.righHight) + 1
					} else {
						B.righHight = 0
					}
					//}
					//if A != nil {
					if A.left != nil {
						A.leftHight = max(A.left.leftHight, A.left.righHight) + 1
					} else {
						A.leftHight = 0
					}
					//}
					C.leftHight = max(C.left.leftHight, C.left.righHight) + 1
					C.righHight = max(C.right.leftHight, C.right.righHight) + 1
					if curRootParent != nil {
						if curRootPosition == PositionLeft {
							curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
						} else {
							curRootParent.righHight = max(curRootParent.right.leftHight, curRootParent.right.righHight) + 1
						}
					}
				}
				return false
			}
			return isChange

		} else {
			curRoot.left = &AVLNode{data: data}
			curRoot.leftHight += 1
			if curRoot.right == nil {
				return true
			}
			return false
		}
	} else {
		if curRoot.right != nil {
			isChange = tree.insert(curRoot, PositionRight, curRoot.right, data)
			curRoot.righHight = max(curRoot.right.leftHight, curRoot.right.righHight) + 1
			if !isChange {
				return
			}
			//curRoot.righHight += 1
			//check and shift
			if abs(curRoot.leftHight, curRoot.righHight) >= 2 {
				insertPosition := PositionInit
				if data <= curRoot.right.data {
					insertPosition = PositionLeft
				} else {
					insertPosition = PositionRight
				}
				//RR
				if insertPosition == PositionRight {
					// A
					//    B
					//       C
					A := curRoot
					B := curRoot.right
					//C := curRoot.right.right

					//step1:将B节点作为中心节点，放到A的位置
					if curRootParent == nil {
						tree.root = B
					} else {
						if curRootPosition == PositionLeft {
							curRootParent.left = B
						} else {
							curRootParent.right = B
						}
					}
					//step2:将A的右孩子作为B的左孩子，B的左孩子作为A的右孩子
					A.right = B.left
					B.left = A
					//step3:调整各自的leftHeight和rightHeight
					if A.right != nil {
						A.righHight = max(A.right.leftHight, A.right.righHight) + 1
					} else {
						A.righHight = 0
					}

					B.leftHight = max(B.left.leftHight, B.left.righHight) + 1
					if curRootParent != nil {
						if curRootPosition == PositionLeft {
							curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
						} else {
							curRootParent.righHight = max(curRootParent.right.leftHight, curRootParent.right.righHight) + 1
						}
					}

				} else {
					//RL

					//  A
					//     B
					//  C
					A := curRoot
					B := curRoot.right
					C := curRoot.right.left
					//step1:将C放到A的位置
					if curRootParent == nil {
						tree.root = C
					} else {
						if curRootPosition == PositionLeft {
							curRootParent.left = C
						} else {
							curRootParent.right = C
						}
					}
					//step2:将A的right指向C的left
					A.right = C.left
					//step3:将B的left指向C的right
					B.left = C.right
					//step4:将A作为C的左孩子，B作为C的右孩子
					C.left = A
					C.right = B
					//step5:调整相关节点的高度
					if A.right != nil {
						A.righHight = max(A.right.leftHight, A.right.righHight) + 1
					} else {
						A.righHight = 0
					}
					if B.left != nil {
						B.leftHight = max(B.left.leftHight, B.left.righHight) + 1
					} else {
						B.leftHight = 0
					}
					C.leftHight = max(C.left.leftHight, C.left.righHight) + 1
					C.righHight = max(C.right.leftHight, C.right.righHight) + 1
					if curRootParent != nil {
						if curRootPosition == PositionLeft {
							curRootParent.leftHight = max(curRootParent.left.leftHight, curRootParent.left.righHight) + 1
						} else {
							curRootParent.righHight = max(curRootParent.leftHight, curRootParent.righHight) + 1
						}
					}
				}
				return false
			}
			return isChange
		} else {
			curRoot.right = &AVLNode{data: data}
			curRoot.righHight += 1
			if curRoot.left == nil {
				return true
			}
			return false
		}
	}

}

func VLR(root *AVLNode) {
	if root == nil {
		return
	}
	fmt.Println(root.data)
	VLR(root.left)
	VLR(root.right)
}
func LDR(root *AVLNode) {
	if root == nil {
		return
	}
	LDR(root.left)
	fmt.Println(root.data)
	LDR(root.right)

}
func PTR(root *AVLNode) {
	if root == nil {
		return
	}
	fmt.Println(root.data)
	PTR(root.left)
	PTR(root.right)

}
func abs(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
