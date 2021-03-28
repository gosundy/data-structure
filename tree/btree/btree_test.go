package btree

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"strings"
	"testing"

	"github.com/awalterschulze/gographviz"
)

func TestNewBtree(t *testing.T) {
	tree := NewBtree(5)
	for i := 0; i < 10; i++ {
		if i == 6 {
			tree.Insert(i)
		} else {
			tree.Insert(i)
		}
	}
	generateGraphPic(tree, "tmp.png")
}
func TestNewBtree02(t *testing.T) {
	tree := NewBtree(4)
	count := 40
	datas := rand.Perm(count)
	for i := 0; i < count; i++ {
		tree.Insert(datas[i])
		generateGraphPic(tree, fmt.Sprintf("tmp%d.png", i))
	}

}
func TestBTree_Delete(t *testing.T) {
	tree := NewBtree(4)
	count := 10
	datas := rand.Perm(count)
	for i := 0; i < count; i++ {
		tree.Insert(datas[i])
	}
	generateGraphPic(tree, fmt.Sprintf("tmp.png"))
	for i := 0; i < count; i++ {
		if i==7{
			tree.Delete(i)
		}else{
			tree.Delete(i)
		}

		generateGraphPic(tree, fmt.Sprintf("tmp%d.png", i))
	}
}
func TestBTree_Delete01(t *testing.T) {
	tree := NewBtree(20)
	count := 2000
	datas := rand.Perm(count)
	for i := 0; i < count; i++ {
		tree.Insert(datas[i])
	}
	generateGraphPic(tree, fmt.Sprintf("tmp.svg"))
	for i := 0; i < count; i++ {
	//	t.Logf("idx %d to be deleted %d\n", i, datas[i])
		if i==4{
			tree.Delete(datas[i])
		}else{
			tree.Delete(datas[i])
		}

		//generateGraphPic(tree, fmt.Sprintf("tmp%d.png", i))
		LDR(tree.root)
	}
}
func LDR(node *BTreeNode) {
	if node == nil {
		return
	}
	for i := 0; i < node.keyLen; i++ {
		LDR(node.childs[i])
		fmt.Print(node.keys[i], " ")
	}
	LDR(node.childs[node.keyLen])
}
func generateGraphPic(tree *BTree, picName string) {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	generateGraphPicVLR(tree.root, graph)
	//fix
	buffer := bytes.Buffer{}
	buffer.WriteString(`digraph G {
    node [shape = record,height=.1];
`)
	buffer.WriteString(string(graph.String())[len("digraph G {")+1:])
	// 输出文件
	ioutil.WriteFile("tmp.gv", buffer.Bytes(), 0666)

	// 产生图片
	ss:=strings.Split(picName,".")
	if ss[len(ss)-1]=="png"{
		system(fmt.Sprintf("dot tmp.gv -T png -o %s", picName))
	}
	if ss[len(ss)-1]=="svg"{
		system(fmt.Sprintf("dot tmp.gv -T svg -o %s", picName))
	}
	system(fmt.Sprintf("open %s", picName))
}
func generateGraphPicVLR(root *BTreeNode, graph *gographviz.Graph) {
	if root == nil {
		return
	}
	nodeName := graphGetNode(root)
	nodeLabel := graphGetNodeLabel(root)
	graph.AddNode("", nodeName, map[string]string{"label": nodeLabel})
	for i := 0; i <= root.keyLen; i++ {
		if root.childs[i] != nil {
			_nodeName := graphGetNode(root.childs[i])
			_nodeLabel := graphGetNodeLabel(root.childs[i])
			graph.AddNode("", _nodeName, map[string]string{"label": _nodeLabel})
			graph.AddPortEdge(nodeName, fmt.Sprintf("f%d", i), _nodeName, "", true, nil)
			generateGraphPicVLR(root.childs[i], graph)
		}

	}

}
func graphGetNode(node *BTreeNode) string {
	s := "\""
	for i := 0; i < node.keyLen-1; i++ {
		s = s + fmt.Sprintf("%d", node.keys[i]) + "_"
	}
	s = s + fmt.Sprintf("%d", node.keys[node.keyLen-1]) + "\""
	return s
}
func graphGetNodeLabel(node *BTreeNode) string {
	s := "\""
	for i := 0; i < node.keyLen; i++ {
		s = s + fmt.Sprintf("<f%d> |%d|", i, node.keys[i])
	}
	s = s + fmt.Sprintf(" <f%d>", node.keyLen)
	s = s + "\""
	return s
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
