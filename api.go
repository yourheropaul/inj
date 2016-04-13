// +build !noglobals

package inj

//////////////////////////////////////////////
// Interface definitions
//////////////////////////////////////////////

// A Grapher is anything that can represent an application graph
type Grapher interface {
	Provide(inputs ...interface{}) error
	Inject(fn interface{}, args ...interface{})
	Assert() (valid bool, errors []string)
	AddDatasource(...interface{}) error
}

//////////////////////////////////////////////
// The one true global variable
//////////////////////////////////////////////

// A default grapher to use in the public API
var globalGraph Grapher = NewGraph()

// Fetch the current grapher instance (in other words, get the global graph)
func GetGrapher() Grapher {
	return globalGraph
}

// Set a specific grapher instance, which will replace the global graph.
func SetGrapher(g Grapher) {
	globalGraph = g
}

//////////////////////////////////////////////
// Public API
//////////////////////////////////////////////

// Insert zero or more objected into the graph, and then attempt to wire up any unmet
// dependencies in the graph.
//
// As explained in the main documentation (https://godoc.org/github.com/yourheropaul/inj),
// a graph consists of what is essentially a map of types to values. If the same type is
// provided twice with different values, the *last* value will be stored in the graph.
func Provide(inputs ...interface{}) error {
	return globalGraph.Provide(inputs...)
}

// Given a function, call it with arguments assigned
// from the graph. Additional arguments can be provided
// for the sake of utility.
//
// Inject() will panic if the provided argument isn't a function,
// or if the provided function accepts variadic arguments (because
// that's not currently supported in the scope of inj).
func Inject(fn interface{}, args ...interface{}) {
	globalGraph.Inject(fn, args...)
}

// Make sure that all provided dependencies have their
// requirements met, and return a list of errors if they
// haven't. A graph is never really finalised, so Provide() and
// Assert() can be called any number of times.
func Assert() (valid bool, errors []string) {
	return globalGraph.Assert()
}

// Add any number of Datasources, DatasourceReaders or DatasourceWriters
// to the graph. Returns an error if any of the supplied arguments aren't
// one of the accepted types.
//
// Once added, the datasources will be active immediately, and the graph
// will automatically re-Provide itself, so that any depdendencies that
// can only be met by an external datasource will be wired up automatically.
//
func AddDatasource(ds ...interface{}) error {
	return globalGraph.AddDatasource(ds...)
}
