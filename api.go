package inj

//////////////////////////////////////////////
// Interface definitions
//////////////////////////////////////////////

// A Grapher is anything that can represent an application graph
type Grapher interface {
	Provide(inputs ...interface{}) error
	// Inject
	// Assert
}

//////////////////////////////////////////////
// The one global variable
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
