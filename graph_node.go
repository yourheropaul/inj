package inj

import "reflect"

// Should this be exported?
type GraphNode struct {
	Name         string
	Object       interface{}
	Type         reflect.Type
	Value        reflect.Value
	Dependencies []GraphNodeDependency
}

type nodeMap map[reflect.Type]*GraphNode

func NewGraphNode() (n *GraphNode) {

	n = &GraphNode{}

	n.Dependencies = make([]GraphNodeDependency, 0)

	return
}
