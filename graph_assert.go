package inj

// Make sure that all provided dependencies have their
// requirements met, and return a list of errors if they
// haven't. A graph is never really finalised, so Provide() and
// Assert() can be called any number of times.
func (g *Graph) Assert() (valid bool, errors []string) {

	valid = true

	if g.UnmetDependencies > 0 || len(g.Errors) > 0 {
		valid = false
	}

	return valid, g.Errors
}
