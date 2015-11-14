package inj

import "reflect"

type graphNode struct {
	Name         string
	Object       interface{}
	Type         reflect.Type
	Value        reflect.Value
	Dependencies []graphNodeDependency
}

type nodeMap map[reflect.Type]*graphNode

func newGraphNode() (n *graphNode) {

	n = &graphNode{}

	n.Dependencies = make([]graphNodeDependency, 0)

	return
}
