package inj

import (
	"fmt"
	"reflect"
)

// Given a function, call it with arguments from the graph
func (g *Graph) Inject(fn interface{}, args ...interface{}) {

	// Reflect the input
	f := reflect.ValueOf(fn)

	// We can only accept functions
	if f.Type().Kind() != reflect.Func {
		panic("[inj.Inject] Passed argument is not a function")
	}

	// Variadic functions aren't currently supported
	if f.Type().Kind() != reflect.Func {
		panic("[inj.Inject] Passed function is variadic")
	}

	// Assemble a list of function arguments
	argv := make([]reflect.Value, 0)

	// Assemble extra arg types list
	xargs := make([]reflect.Type, len(args))

	for i := 0; i < len(args); i++ {
		xargs[i] = reflect.TypeOf(args[i])
	}

	// Number of required incoming arguments
	argc := f.Type().NumIn()

	for i := 0; i < argc; i++ {

		found := false

		// Get an incoming arg reflection type
		in := f.Type().In(i)

		// Find an entry in the graph
		for typ, node := range g.Nodes {
			if typ.AssignableTo(in) {
				argv = append(argv, node.Value)
				found = true
				break
			}
		}

		// CHeck the additional args, if available
		if !found && len(xargs) > 0 {

			// Look in the additional args for the requirement
			for i, xarg := range xargs {
				if xarg.AssignableTo(in) {
					argv = append(argv, reflect.ValueOf(args[i]))
					found = true
					break
				}
			}
		}

		// If it's STILL not found, panic
		if !found {
			panic(fmt.Sprintf("[inj.Inject] Can't find value for arg %d [%s]", i, in))
		}
	}

	// Make the function call, with the args which should
	// now be complete.
	f.Call(argv)
}
