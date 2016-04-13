package inj

import (
	"fmt"
	"reflect"
)

// Usually called after Provide() to assign the values
// of all requested dependencies.
func (g *graph) connect() {

	// Reset error counts
	g.unmetDependency = 0
	g.errors = make([]string, 0)

	// loop through all nodes
	for _, node := range g.nodes {

		// assign dependencies to the object
		for _, dep := range node.Dependencies {
			if e := g.assignValueToNode(node.Value, dep); e != nil {
				g.unmetDependency++
				g.errors = append(g.errors, e.Error())
			}
		}
	}
}

func (g *graph) assignValueToNode(o reflect.Value, dep graphNodeDependency) error {

	parents := []reflect.Value{}
	v, err := g.findFieldValue(o, dep.Path, &parents)
	vtype := v.Type()

	if err != nil {
		return err
	}

	// Sanity check
	if !v.CanSet() {
		return fmt.Errorf("%s%s can't be set", o, dep.Path)
	}

	// If there are any datasource paths supplied...
	for _, path := range dep.DatasourcePaths {

		// ...check to see if a datasource reader has the value
		for _, d := range g.datasourceReaders {

			if dsvalue, err := d.Read(path); err == nil {

				typ := reflect.TypeOf(dsvalue)

				value := reflect.ValueOf(dsvalue)

				if typ != vtype && typ.ConvertibleTo(vtype) {
					value = value.Convert(vtype)
				}

				if value.Type().AssignableTo(vtype) {

					// The value can be set by reflection
					v.Set(value)

					// Any datasourcewriters need to be updated
					for _, w := range g.datasourceWriters {
						w.Write(path, v.Interface())
					}

					return nil
				}
			}
		}
	}

	// Run through the graph and see if anything is settable
	for typ, node := range g.nodes {

		valid := true

		// Don't assign anything to itself or its children
		for _, parent := range parents {

			if parent.Interface() == node.Value.Interface() {
				valid = false
				break
			}
		}

		if !valid {
			continue
		}

		if typ.AssignableTo(v.Type()) {

			// The value can be set by reflection
			v.Set(node.Value)

			// Any datasourcewriters need to be updated
			for _, path := range dep.DatasourcePaths {
				for _, w := range g.datasourceWriters {
					w.Write(path, v.Interface())
				}
			}

			return nil
		}
	}

	return fmt.Errorf("Couldn't find suitable dependency for %s", dep.Type)
}

// Required a struct type
func (g *graph) findFieldValue(parent reflect.Value, path structPath, linneage *[]reflect.Value) (reflect.Value, error) {

	*linneage = append(*linneage, parent)

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
	return g.findFieldValue(f, path, linneage)
}
