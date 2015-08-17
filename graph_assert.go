package inj

func (g *Graph) Assert() (valid bool, errors []string) {

	valid = true

	if g.UnmetDependencies > 0 || len(g.Errors) > 0 {
		valid = false
	}

	return valid, g.Errors
}
