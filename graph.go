package inj

import "reflect"

// A Graph object represents an flat tree of application
// dependencies, a count of currently unmet dependencies,
// and a list of encountered errors.
type graph struct {
	nodes             nodeMap
	unmetDependency   int
	errors            []string
	indexes           []reflect.Type
	datasourceReaders []DatasourceReader
	datasourceWriters []DatasourceWriter
}

// Create a new instance of a graph with allocated memory
func NewGraph(providers ...interface{}) Grapher {

	g := &graph{}

	g.nodes = make(nodeMap)
	g.errors = make([]string, 0)
	g.datasourceReaders = make([]DatasourceReader, 0)
	g.datasourceWriters = make([]DatasourceWriter, 0)

	g.Provide(providers...)

	return g
}

// Add a node by reflection type
func (g *graph) add(typ reflect.Type) (n *graphNode) {

	n = newGraphNode()
	g.nodes[typ] = n

	return
}
