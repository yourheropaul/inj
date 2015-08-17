package inj

import "testing"

// New graph nodes should be different objects
func Test_GraphNodeInitialisation1(t *testing.T) {

	gn1, gn2 := NewGraphNode(), NewGraphNode()

	if gn1 == gn2 {
		t.Errorf("gn1 == gn2")
	}
}
