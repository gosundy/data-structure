package graph

import "testing"

func TestNewGraph(t *testing.T) {
	g:=NewGraph()
	g.AddEdge("a","b",1).
		AddEdge("b","a",2).
		AddEdge("a","c",3).
		AddEdge("c","b",4).Visit()
}
