package redtree

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"testing"

	"github.com/awalterschulze/gographviz"
)

func TestRedBlackTree_Insert(t *testing.T) {
	datas := rand.Perm(10)
	rbTree := &RedBlackTree{}
	for i := 0; i < len(datas); i++ {
		rbTree.Insert(datas[i])
		generateGraphPic(rbTree, fmt.Sprintf("tmp-%d.png", i))
	}

}
func TestRedBlackTree_Insert01(t *testing.T) {
	datas := rand.Perm(100)
	rbTree := &RedBlackTree{}
	for i := 0; i < len(datas); i++ {
		if i == 21 {
			rbTree.Insert(datas[i])
		} else {
			rbTree.Insert(datas[i])
		}

		generateGraphPic(rbTree, fmt.Sprintf("tmp-%d.png", i))
	}

}
func TestRedBlackTree_Delete(t *testing.T) {
	datas := rand.Perm(40)
	rbTree := &RedBlackTree{}
	for i := 0; i < len(datas); i++ {
		rbTree.Insert(datas[i])
	}
	generateGraphPic(rbTree, fmt.Sprintf("tmp.png"))
	for i := 0; i < len(datas); i++ {
		if i == 12 {
			rbTree.Delete(i)
		} else {
			rbTree.Delete(i)
		}
		generateGraphPic(rbTree, fmt.Sprintf("tmp-%d.png", i))
	}
}
func TestRedBlackTree_Delete2(t *testing.T) {
	graphAst, _ := gographviz.ParseString(`digraph G {node [shape = record,height=.1];}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}

	parentNode := "node01"
	node02 := "node02"
	graph.AddPortEdge(parentNode,"f0",node02,"",true,nil)
	graph.AddNode("", parentNode, map[string]string{"label": `"<f0> |G|<f1> |M|<f2> |P|<f3> |X|<f4>"`})
	graph.AddNode("", node02, map[string]string{"label": `"<f0> |A|<f1> |C|<f2> |D|<f3> |E|<f4>"`})

	graph.AddNode("","node", map[string]string{"shape":"record,height=.1"})
	// 输出文件
	ioutil.WriteFile("tmp2.gv", []byte(graph.String()), 0666)

	// 产生图片
	system(fmt.Sprintf("dot tmp2.gv -T png -o %s", "a.png"))
	system(fmt.Sprintf("open %s", "a.png"))
}
func generateGraphPic(tree *RedBlackTree, picName string) {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	generateGraphPicVLR(tree.root, graph)
	// 输出文件
	ioutil.WriteFile("tmp.gv", []byte(graph.String()), 0666)

	// 产生图片
	system(fmt.Sprintf("dot tmp.gv -T png -o %s", picName))
	system(fmt.Sprintf("open %s", picName))
}
func generateGraphPicVLR(root *TreeNode, graph *gographviz.Graph) {
	if root == nil {
		return
	}
	attrCircleColorRed := map[string]string{"shape": "circle", "color": "red"}
	attrCircleColorBlack := map[string]string{"shape": "circle", "color": "black"}
	if root.color == red {
		graph.AddNode(graph.Name, fmt.Sprintf("%v", root.data), attrCircleColorRed)
	} else {
		graph.AddNode(graph.Name, fmt.Sprintf("%v", root.data), attrCircleColorBlack)
	}
	if root.left != nil {
		if root.left.color == red {
			graph.AddNode(graph.Name, fmt.Sprintf("%v", root.left.data), attrCircleColorRed)
		} else {
			graph.AddNode(graph.Name, fmt.Sprintf("%v", root.left.data), attrCircleColorBlack)
		}

		graph.AddEdge(fmt.Sprintf("%v", root.data), fmt.Sprintf("%v", root.left.data), true, nil)
	}
	if root.right != nil {
		if root.right.color == red {
			graph.AddNode(graph.Name, fmt.Sprintf("%v", root.right.data), attrCircleColorRed)
		} else {
			graph.AddNode(graph.Name, fmt.Sprintf("%v", root.right.data), attrCircleColorBlack)
		}

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
