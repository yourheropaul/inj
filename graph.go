package inj

import "reflect"

// A Graph object represents an flat tree of application
// dependencies, a count of currently unmet dependencies,
// and a list of encountered errors.
type Graph struct {
	Nodes             nodeMap
	UnmetDependencies int
	Errors            []string
	indexes           []reflect.Type
}

// Create a new instance of a graph with allocated memory
func NewGraph(providers ...interface{}) (g *Graph) {

	g = &Graph{}

	g.Nodes = make(nodeMap)
	g.Errors = make([]string, 0)

	g.Provide(providers...)

	return
}

// Add a node by reflection type
func (g *Graph) add(typ reflect.Type) (n *GraphNode) {

	n = NewGraphNode()
	g.Nodes[typ] = n

	return
}
