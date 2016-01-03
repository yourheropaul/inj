package inj

import "reflect"

// Insert zero or more objected into the graph, and then attempt to wire up any unmet
// dependencies in the graph.
//
// As explained in the main documentation (https://godoc.org/github.com/yourheropaul/inj),
// a graph consists of what is essentially a map of types to values. If the same type is
// provided twice with different values, the *last* value will be stored in the graph.
func (g *Graph) Provide(inputs ...interface{}) error {

	for _, input := range inputs {

		// Get reflection types
		mtype, stype := getReflectionTypes(input)

		// Assign a node in the graph
		n := g.add(mtype)

		// Populate the new node
		n.Object = input
		n.Type = mtype
		n.Value = reflect.ValueOf(input)
		n.Name = identifier(stype)

		// For structs, find dependencies
		if stype.Kind() == reflect.Struct {
			var basePath = emptyStructPath()
			findDependencies(stype, &n.Dependencies, &basePath)
		}
	}

	///////////////////////////////////////////////
	// Plug everything together
	///////////////////////////////////////////////

	g.connect()

	///////////////////////////////////////////////
	// Store a list of types for speed later on
	///////////////////////////////////////////////

	g.indexes = make([]reflect.Type, 0)

	for typ, _ := range g.nodes {
		g.indexes = append(g.indexes, typ)
	}

	return nil
}
