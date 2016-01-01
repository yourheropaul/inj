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
var graph Grapher = NewGraph()

// Fetch the current grapher instance
func GetGrapher() Grapher {
	return graph
}

// Set a specific grapher instance
func SetGrapher(g Grapher) {
	graph = g
}

//////////////////////////////////////////////
// Public API
//////////////////////////////////////////////

// Insert a set of arbitrary objects into the
// application graph
func Provide(inputs ...interface{}) error {
	return graph.Provide(inputs...)
}

// Given a function, call it with arguments assigned
// from the graph. Additional arguments can be provided
// for the sake of utility.
func Inject(fn interface{}, args ...interface{}) {
	graph.Inject(fn, args...)
}

// Make sure that all provided dependencies have their
// requirements met, and return a list of errors if they
// don't.
func Assert() (valid bool, errors []string) {
	return graph.Assert()
}

// Add zero or more datasources to the global graph
func AddDatasource(ds ...interface{}) error {
	return graph.AddDatasource(ds)
}
