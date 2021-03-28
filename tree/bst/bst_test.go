package bst

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBst_Insert(t *testing.T) {
	bst := Bst{}
	count := 100
	datas := rand.Perm(count)
	for i := 0; i < count; i++ {
		bst.Insert(datas[i])
	}
	LDR(bst.root)
}
func TestBst_Delete(t *testing.T) {
	bst := Bst{}
	count := 100
	datas := rand.Perm(count)
	for i := 0; i < count; i++ {
		bst.Insert(datas[i])
	}
	LDR(bst.root)
	fmt.Println("\ninsert completed")
	for i := 0; i < count; i++ {
		bst.Delete(i)
		LDR(bst.root)
		fmt.Println("")
	}

}
