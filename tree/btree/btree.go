package btree

type Btree struct {
	root *BTreeNode
}
type BTreeNode struct {
	parents []*BTreeNode
	childs  []*BTreeNode
}
