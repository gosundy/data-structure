package avl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os/exec"
	"testing"

	"github.com/awalterschulze/gographviz"
)

func TestAVLTree(t *testing.T) {
	tree := NewAvlTree()
	datas := rand.Perm(10)
	for _, data := range datas {
		tree.Put(data)
	}
	LDR(tree.root)
}
func TestAVLTree_Put01(t *testing.T) {
	t.Log("LL")
	tree := NewAvlTree()
	count:=int(math.Exp2(8))-1
	for i := 0; i < count; i++ {
		tree.Put(i)
	}
	generateGraphPic(tree)
}
func TestAVLTree_Put02(t *testing.T) {
	t.Log("RR")
	tree := NewAvlTree()
	for i := 100; i > 0; i-- {
		tree.Put(i)
	}
	generateGraphPic(tree)
}
func TestAVLTree_Put03(t *testing.T) {
	t.Log("RAND")
	tree := NewAvlTree()
	//datas:=[]int{12,4,2,13,10,0,3,11,7,5,15,1,9,14,6}
	datas:=[]int{12,4,2,13,10,0,3,11,7,5,15,1,9,14,6}
	count := len(datas)
	//datas := rand.Perm(count)
	for i := 0; i < count; i++ {

		//if i==15{
		//	tree.Put(datas[i])
		//}
		tree.Put(datas[i])
	}
	t.Log(datas)
	generateGraphPic(tree)
}
func TestAVLTree_Put04(t *testing.T) {
	t.Log("RAND")
	tree := NewAvlTree()
	//datas:=[]int{12,4,2,13,10,0,3,11,7,5,15,1,9,14,6}
	//datas:=[]int{12,4,2,13,10,0,3,11,7,5,15,1,9,14,6}

	datas := rand.Perm(1000)
	count := len(datas)
	for i := 0; i < count; i++ {

		//if i==15{
		//	tree.Put(datas[i])
		//}
		tree.Put(datas[i])
	}
	t.Log(datas)
	generateGraphPic(tree)
}
func generateGraphPic(tree *AVLTree) {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	generateGraphPicVLR(tree.root, graph)
	// 输出文件
	ioutil.WriteFile("tmp.gv", []byte(graph.String()), 0666)

	// 产生图片
	system("dot tmp.gv -T png -o tmp.png")
	system("open tmp.png")
}
func generateGraphPicVLR(root *AVLNode, graph *gographviz.Graph) {
	if root == nil {
		return
	}
	attrCircle := map[string]string{"shape": "circle"}
	graph.AddNode(graph.Name, fmt.Sprintf("%v", root.data), attrCircle)
	if root.left != nil {
		graph.AddNode(graph.Name, fmt.Sprintf("%v", root.left.data), attrCircle)
		graph.AddEdge(fmt.Sprintf("%v", root.data), fmt.Sprintf("%v", root.left.data), true, nil)
	}
	if root.right != nil {
		graph.AddNode(graph.Name, fmt.Sprintf("%v", root.right.data), attrCircle)
		graph.AddEdge(fmt.Sprintf("%v", root.data), fmt.Sprintf("%v", root.right.data), true, nil)
	}
	generateGraphPicVLR(root.left, graph)
	generateGraphPicVLR(root.right, graph)

}

//调用系统指令的方法，参数s 就是调用的shell命令
func system(s string) {
	cmd := exec.Command(`/bin/sh`, `-c`, s) //调用Command函数
	var out bytes.Buffer                    //缓冲字节

	cmd.Stdout = &out //标准输出
	err := cmd.Run()  //运行指令 ，做判断
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", out.String()) //输出执行结果
}
