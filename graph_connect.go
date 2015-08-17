package inj

import (
	"fmt"
	"reflect"
)

// Usually called after Provide() to assign the values
// of all requested dependencies.
func (g *Graph) connect() {

	// Reset error counts
	g.UnmetDepdendencies = 0
	g.Errors = make([]string, 0)

	// loop through all nodes
	for _, node := range g.Nodes {

		// assign depdendencies to the object
		for _, dep := range node.Dependencies {
			if e := g.assignValueToNode(node.Value, dep); e != nil {
				g.UnmetDepdendencies++
				g.Errors = append(g.Errors, e.Error())
			}
		}
	}
}

func (g *Graph) assignValueToNode(o reflect.Value, dep GraphNodeDependency) error {

	v, err := g.findFieldValue(o, dep.Path)

	if err != nil {
		return err
	}

	// If value has already been set, then skip it
	if !zero(v) {
		return nil
	}

	// Sanity check
	if !v.CanSet() {
		return fmt.Errorf("%s%s can't be set", o, dep.Path)
	}

	// Run through the graph and see if anything is settable
	for typ, node := range g.Nodes {

		if typ.AssignableTo(v.Type()) {
			v.Set(node.Value)
			return nil
		}
	}

	return fmt.Errorf("Couldn't find suitable depdendency for %s", dep.Type)
}

// Required a struct type
func (g *Graph) findFieldValue(parent reflect.Value, path StructPath) (reflect.Value, error) {

	// Dereference incoming values
	if parent.Kind() == reflect.Ptr {
		parent = parent.Elem()
	}

	// Only accept structs
	if parent.Kind() != reflect.Struct {
		return parent, fmt.Errorf("Type is %s, not struct", parent.Kind().String())
	}

	// Take the first entry from the path
	stub, path := path.Shift()

	// Try to get the field
	f := parent.FieldByName(stub)

	if !f.IsValid() {
		return f, fmt.Errorf("Can't find field %s in %s", stub, parent)
	}

	// If that's the end of the path, return the value
	if path.Empty() {
		return f, nil
	}

	// Otherwise recurse
	return g.findFieldValue(f, path)
}
