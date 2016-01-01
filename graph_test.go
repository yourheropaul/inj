package inj

import (
	"reflect"
	"testing"
)

func assertNoGraphErrors(t *testing.T, g *Graph) {

	if len(g.Errors) != 0 {
		t.Fatalf("Graph was initialised with errors > 0")
	}

	if g.UnmetDependencies != 0 {
		t.Fatalf("Graph was initialised with UnmetDependencies > 0")
	}
}

// New graphs should be different objects
func Test_GraphInitialisation1(t *testing.T) {

	g1, g2 := NewGraph(), NewGraph()

	if g1 == g2 {
		t.Errorf("g1 == g2")
	}
}

// Initial graph state should be zero values
func Test_GraphInitialisation2(t *testing.T) {

	g := NewGraph()

	if len(g.Errors) != 0 {
		t.Errorf("Graph was initialised with errors > 0")
	}

	if g.UnmetDependencies != 0 {
		t.Errorf("Graph was initialised with UnmetDependencies > 0")
	}
}

// Added nodes should be in the graph object
func Test_GraphNodeAddition(t *testing.T) {

	g := NewGraph()
	typ := reflect.TypeOf(0)

	if len(g.nodes) != 0 {
		t.Errorf("Graph was initialised with nodes > 0")
	}

	n := g.add(typ)

	if n == nil {
		t.Errorf("New graph node is nil")
	}

	if len(g.nodes) != 1 {
		t.Errorf("New graph node count != 1")
	}

	for typ2, node := range g.nodes {

		if typ2 == typ && node != n {
			t.Errorf("Expected typ to equate to node")
		}

		if node == n && typ2 != typ {
			t.Errorf("Expected node to equate to typ")
		}
	}
}
