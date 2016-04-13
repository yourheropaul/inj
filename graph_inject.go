package inj

import (
	"fmt"
	"reflect"
)

// Given a function, call it with arguments from the graph.
// Throws a runtime error in the form of a panic on failure.
func (g *graph) Inject(fn interface{}, args ...interface{}) {

	// Reflect the input
	f := reflect.ValueOf(fn)

	// It's slightly faster to store the type rather than constantly
	// retrieving it.
	ftype := f.Type()

	// We can only accept functions
	if ftype.Kind() != reflect.Func {
		panic("[inj.Inject] Passed argument is not a function")
	}

	// Variadic functions aren't currently supported
	if ftype.IsVariadic() {
		panic("[inj.Inject] Passed function is variadic")
	}

	// Assemble extra arg types list
	xargs := make([]reflect.Type, len(args))

	for i := 0; i < len(args); i++ {
		xargs[i] = reflect.TypeOf(args[i])
	}

	// Number of required incoming arguments
	argc := ftype.NumIn()

	// Assemble a list of function arguments
	argv := make([]reflect.Value, argc)

	for i := 0; i < argc; i++ {

		func() {
			// Get an incoming arg reflection type
			in := ftype.In(i)

			// Look in the additional args list for the requirement
			for j := 0; j < len(xargs); j++ {
				if xargs[j].AssignableTo(in) {
					argv[i] = reflect.ValueOf(args[j])
					return
				}
			}

			// Find an entry in the graph
			for j := 0; j < len(g.indexes); j++ {
				if g.indexes[j].AssignableTo(in) {
					argv[i] = g.nodes[g.indexes[j]].Value
					return
				}
			}

			// If it's STILL not found, panic
			panic(fmt.Sprintf("[inj.Inject] Can't find value for arg %d [%s]", i, in))
		}()
	}

	// Make the function call, with the args which should now be complete.
	f.Call(argv)
}
