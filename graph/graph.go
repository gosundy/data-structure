package graph

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

var IdxNotFound = errors.New("vertex not found")

type Graph struct {
	VerIndices      map[string]int64
	IdxIndices      map[int64]string
	Vertexes        map[string]struct{}
	Edges           map[int64]map[int64]float64
	curMaxVertexId  int64
	idGeneratorLock sync.RWMutex
	edgesLock       sync.RWMutex
}

func NewGraph() *Graph {
	g := &Graph{}
	g.VerIndices = make(map[string]int64)
	g.IdxIndices = make(map[int64]string)
	g.Edges = make(map[int64]map[int64]float64)
	g.VerIndices = make(map[string]int64)
	return g
}
func (g *Graph) AddVertex(vertexName string) *Graph {
	g.idGeneratorLock.RLock()
	_, ok := g.VerIndices[vertexName]
	if ok {
		g.idGeneratorLock.RUnlock()
		return g
	}
	g.idGeneratorLock.RUnlock()
	g.idGeneratorLock.Lock()
	defer g.idGeneratorLock.Unlock()
	_, ok = g.VerIndices[vertexName]
	if !ok {
		g.curMaxVertexId++
		g.VerIndices[vertexName] = g.curMaxVertexId
		g.IdxIndices[g.curMaxVertexId] = vertexName
		return g
	}
	return g
}

func (g *Graph) GetVertexIdx(vertexName string) (int64, error) {
	g.idGeneratorLock.RLock()
	defer g.idGeneratorLock.RUnlock()
	id, ok := g.VerIndices[vertexName]
	if ok {
		return id, nil
	}
	return 0, IdxNotFound
}

func (g *Graph) AddEdge(tailVertex string, headVerTex string, weight float64) *Graph {
	g.AddVertex(tailVertex).AddVertex(headVerTex)
	tailIdx, _ := g.GetVertexIdx(tailVertex)
	headIdx, _ := g.GetVertexIdx(headVerTex)
	g.edgesLock.Lock()
	defer g.edgesLock.Unlock()
	_, ok := g.Edges[tailIdx]
	if !ok {
		g.Edges[tailIdx] = make(map[int64]float64)
	}
	g.Edges[tailIdx][headIdx] = weight
	return g
}

func (g *Graph) Visit() {
	buf := bytes.Buffer{}
	g.edgesLock.RLock()
	defer g.edgesLock.RUnlock()
	g.idGeneratorLock.RLock()
	defer g.idGeneratorLock.RUnlock()
	for tailIdx, singleEdges := range g.Edges {
		tailName := g.IdxIndices[tailIdx]
		for headIdx, weight := range singleEdges {
			headName := g.IdxIndices[headIdx]
			buf.WriteString(fmt.Sprintf("(%s -> %s, weight:%f)", tailName, headName, weight))
		}
		buf.WriteString("\n")
	}
	fmt.Println(buf.String())
}
