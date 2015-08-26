package inj

import "reflect"

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
			var basePath = EmptyStructPath()
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

	for typ, _ := range g.Nodes {
		g.indexes = append(g.indexes, typ)
	}

	return nil
}
