package inj

import "reflect"

// A Graph object represents an flat tree of application
// dependencies, a count of currently unmet dependencies,
// and a list of encountered errors.
type Graph struct {
	nodes             nodeMap
	UnmetDependencies int
	Errors            []string
	indexes           []reflect.Type
	datasourceReaders []DatasourceReader
	datasourceWriters []DatasourceWriter
}

// Create a new instance of a graph with allocated memory
func NewGraph(providers ...interface{}) (g *Graph) {

	g = &Graph{}

	g.nodes = make(nodeMap)
	g.Errors = make([]string, 0)
	g.datasourceReaders = make([]DatasourceReader, 0)
	g.datasourceWriters = make([]DatasourceWriter, 0)

	g.Provide(providers...)

	return
}

// Add a node by reflection type
func (g *Graph) add(typ reflect.Type) (n *graphNode) {

	n = newGraphNode()
	g.nodes[typ] = n

	return
}
