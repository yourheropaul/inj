package inj

import (
	"fmt"
	"reflect"
)

// Given a function, call it with arguments from the graph.
// Throws a runtime erro,r in the form of a panic, on failure.
func (g *Graph) Inject(fn interface{}, args ...interface{}) {

	// Reflect the input
	f := reflect.ValueOf(fn)

	// We can only accept functions
	if f.Type().Kind() != reflect.Func {
		panic("[inj.Inject] Passed argument is not a function")
	}

	// Variadic functions aren't currently supported
	if f.Type().IsVariadic() {
		panic("[inj.Inject] Passed function is variadic")
	}

	// Assemble extra arg types list
	xargs := make([]reflect.Type, len(args))

	for i := 0; i < len(args); i++ {
		xargs[i] = reflect.TypeOf(args[i])
	}

	// Number of required incoming arguments
	argc := f.Type().NumIn()

	// Assemble a list of function arguments
	argv := make([]reflect.Value, argc)

	for i := 0; i < argc; i++ {

		func() {
			// Get an incoming arg reflection type
			in := f.Type().In(i)

			// Find an entry in the graph
			for typ, node := range g.Nodes {
				if typ.AssignableTo(in) {
					argv[i] = node.Value
					return
				}
			}

			// Check the additional args, if available
			if len(xargs) > 0 {

				// Look in the additional args list for the requirement
				for j, xarg := range xargs {
					if xarg.AssignableTo(in) {
						argv[i] = reflect.ValueOf(args[j])
						return
					}
				}
			}

			// If it's STILL not found, panic
			panic(fmt.Sprintf("[inj.Inject] Can't find value for arg %d [%s]", i, in))
		}()
	}

	// Make the function call, with the args which should now be complete.
	f.Call(argv)
}
