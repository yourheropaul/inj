package inj

import "reflect"

type Graph struct {
	Nodes              nodeMap
	UnmetDepdendencies int
	Errors             []string
}

func NewGraph() (g *Graph) {

	g = &Graph{}

	g.Nodes = make(nodeMap)
	g.Errors = make([]string, 0)

	return
}

func (g *Graph) add(typ reflect.Type) (n *GraphNode) {

	n = NewGraphNode()
	g.Nodes[typ] = n

	return
}
